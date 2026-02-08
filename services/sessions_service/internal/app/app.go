package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	eventspostgres "github.com/fedotovmax/kafka-lib/adapters/db/postgres/events_postgres"
	"github.com/fedotovmax/kafka-lib/kafka"
	"github.com/fedotovmax/kafka-lib/outbox"
	outboxsender "github.com/fedotovmax/kafka-lib/outbox_sender"
	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db/postgres"
	eventspublisher "github.com/fedotovmax/microservices-shop/sessions_service/internal/events_publisher"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/queries"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/usecases"
	"github.com/medama-io/go-useragent"

	securitypostgres "github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db/postgres/security_postgres"
	sessionspostgres "github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db/postgres/sessions_postgres"
	userspostgres "github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db/postgres/users_postgres"
	grpcadapter "github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/grpc"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/config"
	grpccontroller "github.com/fedotovmax/microservices-shop/sessions_service/internal/controller/grpc_controller"
	kafkacontroller "github.com/fedotovmax/microservices-shop/sessions_service/internal/controller/kafka_controller"
	"github.com/fedotovmax/microservices-shop/sessions_service/pkg/logger"
	"github.com/fedotovmax/pgxtx"
)

type App struct {
	c             *config.AppConfig
	log           *slog.Logger
	postgres      postgres.PostgresPool
	grpc          *grpcadapter.Server
	event         *outbox.Outbox
	producer      kafka.Producer
	consumerGroup kafka.ConsumerGroup
}

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

	sessionsPostgres := sessionspostgres.New(ex, log)
	usersPostgres := userspostgres.New(ex, log)
	securityPostgres := securitypostgres.New(ex, log)

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

	eventSender := outboxsender.New(eventsPostgres, txManager)

	usersQuery := queries.NewUser(usersPostgres)

	sessionsQuery := queries.NewSession(sessionsPostgres)

	trustTokenQuery := queries.NewTrustToken(securityPostgres)

	publisher := eventspublisher.New(eventSender)

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

	sessionController := grpccontroller.NewSession(
		log,
		createSessionUsecase,
		refreshSessionUsecase,
	)

	kafkaConsumerController := kafkacontroller.New(log, createUserUsecase)

	eventProcessor, err := outbox.New(log, producer, eventSender, &outboxConfig)

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

	grpcServer := grpcadapter.New(
		grpcadapter.Config{
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

func (a *App) Run(cancel context.CancelFunc) {

	const op = "app.Run"

	log := a.log.With(slog.String("op", op))

	a.event.Start()
	log.Info("event processor starting")
	a.consumerGroup.Start()
	log.Info("consumer group starting")

	log.Info("Try to start GRPC server:", slog.String("port", fmt.Sprintf("%d", a.c.Port)))

	go func() {
		if err := a.grpc.Start(); err != nil {
			log.Error("Cannot start grpc server", logger.Err(err))
			log.Error("Signal to shutdown")
			cancel()
			return
		}
	}()
}

func (a *App) Stop(ctx context.Context) {

	const op = "app.Stop"

	log := a.log.With(slog.String("op", op))

	lifesycleServices := []*service{
		newService("event processor", a.event),
		newService("kafka producer", a.producer),
		newService("kafka consumer-group", a.consumerGroup),
		newService("postgres", a.postgres),
	}

	stopErrorsChan := make(chan error, len(lifesycleServices))

	err := a.grpc.Stop(ctx)

	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("grpc server stopped successfully!")
	}

	wg := &sync.WaitGroup{}

	for _, s := range lifesycleServices {
		wg.Go(func() {
			err := s.Stop(ctx)
			if err != nil {
				stopErrorsChan <- err
			} else {
				log.Info("successfully stopped:", slog.String("service", s.Name()))
			}
		})
	}

	go func() {
		wg.Wait()
		close(stopErrorsChan)
	}()

	stopErrors := make([]error, 0, len(lifesycleServices))

	for err := range stopErrorsChan {
		stopErrors = append(stopErrors, err)
	}

	if len(stopErrors) == 0 {
		log.Info("All resources are closed successfully, exit app")
	} else {
		log.Error("resource with errors:", slog.Int("stop_errors", len(stopErrors)))
		for _, err := range stopErrors {
			log.Error(err.Error())
		}
	}
}
