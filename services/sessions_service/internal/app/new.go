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
	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db/postgres"
	securityPSQL "github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db/postgres/security"
	sessionsPSQL "github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db/postgres/sessions"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/publisher"

	usersPSQL "github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db/postgres/users"
	grpcAdapter "github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/grpc"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/config"
	grpcController "github.com/fedotovmax/microservices-shop/sessions_service/internal/controllers/grpc"
	kafkaController "github.com/fedotovmax/microservices-shop/sessions_service/internal/controllers/kafka"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/queries"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/usecases"
	"github.com/fedotovmax/pgxtx"
	"github.com/medama-io/go-useragent"
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

	sessionsPostgres := sessionsPSQL.New(ex, log)
	usersPostgres := usersPSQL.New(ex, log)
	securityPostgres := securityPSQL.New(ex, log)

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

	usersQuery := queries.NewUser(usersPostgres)

	sessionsQuery := queries.NewSession(sessionsPostgres)

	trustTokenQuery := queries.NewTrustToken(securityPostgres)

	publisher := publisher.New(eventCreator)

	securityCfg := &usecases.SecurityConfig{
		BlacklistCodeExpDuration: c.BlacklistCodeExpDuration,
		LoginBypassExpDuration:   c.LoginBypassExpDuration,
		BlacklistCodeLength:      c.BlacklistCodeLength,
		LoginBypassCodeLength:    c.LoginBypassCodeLength,
		//TODO: Change this values!!!
		DeviceTrustTokenExpDuration: time.Hour,
		DeviceTrustTokenThreshold:   time.Hour,
	}

	tokensCfg := &usecases.TokenConfig{
		TokenIssuer:            c.TokenIssuer,
		TokenSecret:            c.AccessTokenSecret,
		RefreshExpiresDuration: c.RefreshTokenExpDuration,
		AccessExpiresDuration:  c.AccessTokenExpDuration,
	}

	uaParser := useragent.NewParser()

	addLoginBypassUsecase := usecases.NewAddLoginBypassUsecase(
		log,
		securityCfg,
		securityPostgres,
		publisher,
	)

	addToBlacklistUsecase := usecases.NewAddToBlacklistUsecase(
		log,
		securityCfg,
		securityPostgres,
		publisher,
	)

	checkBypassUsecase := usecases.NewCheckBypassUsecase(
		log,
		securityPostgres,
		addLoginBypassUsecase,
	)

	isSessionRevokedUsecase := usecases.NewIsSessionRevokedUsecase(
		log,
		addToBlacklistUsecase,
	)

	isUserInBlacklistUsecase := usecases.NewIsUserInBlacklistUsecase(
		log,
		addToBlacklistUsecase,
	)

	//TODO:
	//revokeSessionsUsecase := usecases.NewRevokeSessionsUsecase(log, sessionsPostgres)

	createUserUsecase := usecases.NewCreateUserUsecase(
		log, usersPostgres, usersQuery,
	)

	refreshSessionUsecase := usecases.NewRefreshSessionUsecase(
		log,
		tokensCfg,
		isUserInBlacklistUsecase,
		isSessionRevokedUsecase,
		uaParser,
		sessionsPostgres,
		sessionsQuery,
	)

	checkAllSecurityMethodsUsecase := usecases.NewCheckAllSecurityMethodsUsecase(
		log,
		securityCfg,
		isSessionRevokedUsecase,
		isUserInBlacklistUsecase,
		checkBypassUsecase,
		addLoginBypassUsecase,
		trustTokenQuery,
	)

	createSessionUsecase := usecases.NewCreateSessionUsecase(
		log,
		tokensCfg,
		txManager,
		uaParser,
		checkAllSecurityMethodsUsecase,
		sessionsPostgres,
		securityPostgres,
		usersQuery,
	)

	sessionController := grpcController.NewSession(
		log,
		createSessionUsecase,
		refreshSessionUsecase,
	)

	kafkaConsumerController := kafkaController.New(log, createUserUsecase)

	eventProcessor, err := outbox.New(log, producer, eventCreator, &outboxConfig)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	consumerGroup, err := kafka.NewConsumerGroup(&kafka.ConsumerGroupConfig{
		Brokers:             c.KafkaBrokers,
		Topics:              []string{events.USER_EVENTS},
		GroupID:             "sessions-service-app",
		SleepAfterRebalance: time.Second * 2,
		AutoCommit:          true,
	}, log, kafkaConsumerController)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	grpcServer := grpcAdapter.New(
		grpcAdapter.Config{
			Addr: fmt.Sprintf(":%d", c.Port),
		},
		sessionController,
	)

	return &App{
			c:             c,
			log:           log,
			grpc:          grpcServer,
			postgres:      postgresPool,
			consumerGroup: consumerGroup,
			event:         eventProcessor,
			producer:      producer,
		},
		nil
}
