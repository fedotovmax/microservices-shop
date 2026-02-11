package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	eventsPSQL "github.com/fedotovmax/kafka-lib/adapters/db/postgres/events"
	eventcreator "github.com/fedotovmax/kafka-lib/event_creator"
	"github.com/fedotovmax/kafka-lib/kafka"
	"github.com/fedotovmax/kafka-lib/outbox"
	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters/db/postgres"
	emailverificationPSQL "github.com/fedotovmax/microservices-shop/users_service/internal/adapters/db/postgres/email_verification"
	usersPSQL "github.com/fedotovmax/microservices-shop/users_service/internal/adapters/db/postgres/users"
	"github.com/fedotovmax/microservices-shop/users_service/internal/publisher"

	grpcAdapter "github.com/fedotovmax/microservices-shop/users_service/internal/adapters/grpc"
	"github.com/fedotovmax/microservices-shop/users_service/internal/config"

	grpcController "github.com/fedotovmax/microservices-shop/users_service/internal/controllers/grpc"
	kafkaController "github.com/fedotovmax/microservices-shop/users_service/internal/controllers/kafka"
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

	usersPostgres := usersPSQL.New(ex, log)

	emailVerificationPostgres := emailverificationPSQL.New(ex, log)

	eventsPostgres := eventsPSQL.New(ex, log)

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

	eventCreator := eventcreator.New(eventsPostgres, txManager)

	publisher := publisher.New(eventCreator)

	usersQuery := queries.NewUsers(usersPostgres)
	emailVerifyLinkQuery := queries.NewEmailVerifyLink(emailVerificationPostgres)

	emailVerificationConfig := &usecases.EmailConfig{
		EmailVerifyLinkExpiresDuration: c.EmailVerifyLinkExpiresDuration,
	}

	createUserUsecase := usecases.NewCreateUserUsecase(
		txManager,
		log,
		emailVerificationConfig,
		usersPostgres,
		emailVerificationPostgres,
		publisher,
		usersQuery,
	)

	sendNewVerifyEmailLinkUsecase := usecases.NewSendNewEmailVerifyLinkUsecase(
		txManager,
		log,
		emailVerificationConfig,
		usersPostgres,
		emailVerificationPostgres,
		publisher,
		usersQuery,
	)

	sessionActionUsecase := usecases.NewSessionActionUsecase(
		txManager,
		log,
		usersPostgres,
		usersQuery,
	)

	updateProfileUsecase := usecases.NewUpdateProfileUsecase(
		txManager,
		log,
		usersPostgres,
		emailVerificationPostgres,
		publisher,
		usersQuery,
	)

	verifyEmailUsecase := usecases.NewVerifyEmailUsecase(
		txManager,
		log,
		usersPostgres,
		emailVerificationPostgres,
		emailVerifyLinkQuery,
	)

	eventProcessor, err := outbox.New(log, producer, eventCreator, &outboxConfig)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	profileController := grpcController.NewProfile(log, updateProfileUsecase, createUserUsecase, usersQuery)

	sessionController := grpcController.NewSession(log, sessionActionUsecase)

	verificationController := grpcController.NewVerification(log, verifyEmailUsecase, sendNewVerifyEmailLinkUsecase)

	kafkaConsumerController := kafkaController.New(log, &ku{})

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

	grpcServer := grpcAdapter.New(
		grpcAdapter.Config{
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
