package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	eventspostgres "github.com/fedotovmax/kafka-lib/adapters/db/postgres/events_postgres"
	"github.com/fedotovmax/kafka-lib/kafka"
	"github.com/fedotovmax/kafka-lib/outbox"
	outboxsender "github.com/fedotovmax/kafka-lib/outbox_sender"
	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters/db/postgres"
	emailverifylinkpostgres "github.com/fedotovmax/microservices-shop/users_service/internal/adapters/db/postgres/email_verify_link_postgres"
	userspostgres "github.com/fedotovmax/microservices-shop/users_service/internal/adapters/db/postgres/users_postgres"
	grpcadapter "github.com/fedotovmax/microservices-shop/users_service/internal/adapters/grpc"
	"github.com/fedotovmax/microservices-shop/users_service/internal/config"
	grpccontroller "github.com/fedotovmax/microservices-shop/users_service/internal/controller/grpc_controller"
	kafkacontroller "github.com/fedotovmax/microservices-shop/users_service/internal/controller/kafka_controller"
	eventspublisher "github.com/fedotovmax/microservices-shop/users_service/internal/events_publisher"
	"github.com/fedotovmax/microservices-shop/users_service/internal/queries"
	"github.com/fedotovmax/microservices-shop/users_service/internal/usecases"
	"github.com/fedotovmax/pgxtx"
)

func New(c *config.AppConfig, log *slog.Logger) (*App, error) {

	const op = "app.New"

	l := log.With(slog.String("op", op))

	poolConnectCtx, poolConnectCtxCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer poolConnectCtxCancel()

	postgresPool, err := postgres.New(poolConnectCtx, &postgres.Config{
		DSN: c.DBUrl,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	l.Info("Successfully created db pool and connected!")

	txManager, err := pgxtx.Init(postgresPool, log.With(slog.String("op", "transaction.manager")))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ex := txManager.GetExtractor()

	usersPostgres := userspostgres.New(ex, log)

	emailVerifyLinkPostgres := emailverifylinkpostgres.New(ex, log)

	eventsPostgres := eventspostgres.New(ex, log)

	outboxConfig := outbox.SmallBatchConfig

	flushConfig := outboxConfig.GetKafkaFlushConfig()

	producer, err := kafka.NewAsyncProducer(kafka.ProducerConfig{
		Brokers:     c.KafkaBrokers,
		MaxMessages: flushConfig.MaxMessages,
		Frequency:   flushConfig.Frequency,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	outboxSender := outboxsender.New(eventsPostgres, txManager)

	publisher := eventspublisher.New(outboxSender)

	usersQuery := queries.NewUsers(usersPostgres)
	emailVerifyLinkQuery := queries.NewEmailVerifyLink(emailVerifyLinkPostgres)

	emailVerificationConfig := &usecases.EmailConfig{
		EmailVerifyLinkExpiresDuration: c.EmailVerifyLinkExpiresDuration,
	}

	createUserUsecase := usecases.NewCreateUserUsecase(
		txManager,
		log,
		emailVerificationConfig,
		usersPostgres,
		emailVerifyLinkPostgres,
		publisher,
		usersQuery,
	)

	sendNewVerifyEmailLinkUsecase := usecases.NewSendNewEmailVerifyLinkUsecase(
		txManager,
		log,
		emailVerificationConfig,
		usersPostgres,
		emailVerifyLinkPostgres,
		publisher,
		usersQuery,
	)

	sessionActionUsecase := usecases.NewSessionActionUsecase(
		txManager,
		log,
		usersPostgres,
		publisher,
		usersQuery,
	)

	updateProfileUsecase := usecases.NewUpdateProfileUsecase(
		txManager,
		log,
		usersPostgres,
		emailVerifyLinkPostgres,
		publisher,
		usersQuery,
	)

	verifyEmailUsecase := usecases.NewVerifyEmailUsecase(
		txManager,
		log,
		usersPostgres,
		emailVerifyLinkPostgres,
		publisher,
		emailVerifyLinkQuery,
	)

	eventProcessor, err := outbox.New(log, producer, outboxSender, &outboxConfig)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	profileController := grpccontroller.NewProfile(log, updateProfileUsecase, createUserUsecase, usersQuery)

	sessionController := grpccontroller.NewSession(log, sessionActionUsecase)

	verificationController := grpccontroller.NewVerification(log, verifyEmailUsecase, sendNewVerifyEmailLinkUsecase)

	kafkaConsumerController := kafkacontroller.New(log, &ku{})

	consumerGroup, err := kafka.NewConsumerGroup(
		&kafka.ConsumerGroupConfig{
			Brokers: c.KafkaBrokers,
			//TODO:change topics for real
			Topics:              []string{"permissions.events"},
			GroupID:             "users-service-app",
			SleepAfterRebalance: time.Second * 2,
			AutoCommit:          true,
		},
		log,
		kafkaConsumerController,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	grpcServer := grpcadapter.New(
		grpcadapter.Config{
			Addr: fmt.Sprintf(":%d", c.Port),
		},
		profileController,
		sessionController,
		verificationController,
	)

	return &App{
			c:             c,
			log:           log,
			grpc:          grpcServer,
			postgres:      postgresPool,
			event:         eventProcessor,
			producer:      producer,
			consumerGroup: consumerGroup,
		},
		nil
}
