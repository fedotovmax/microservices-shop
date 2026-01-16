package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/fedotovmax/kafka-lib/kafka"
	"github.com/fedotovmax/kafka-lib/outbox"
	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter/db/postgres"
	eventspostgres "github.com/fedotovmax/microservices-shop/users_service/internal/adapter/db/postgres/events_postgres"
	userspostgres "github.com/fedotovmax/microservices-shop/users_service/internal/adapter/db/postgres/users_postgres"
	grpcadapter "github.com/fedotovmax/microservices-shop/users_service/internal/adapter/grpc"
	"github.com/fedotovmax/microservices-shop/users_service/internal/config"
	grpccontroller "github.com/fedotovmax/microservices-shop/users_service/internal/controller/grpc_controller"
	kafkacontroller "github.com/fedotovmax/microservices-shop/users_service/internal/controller/kafka_controller"
	"github.com/fedotovmax/microservices-shop/users_service/internal/usecase"
	"github.com/fedotovmax/microservices-shop/users_service/pkg/logger"

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

// TODO: remove fake usecase
type ku struct{}

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

	storage := usecase.CreateStorage(eventsPostgres, usersPostgres)

	usecases := usecase.NewUsecases(storage, txManager, log, &usecase.Config{
		EmailVerifyLinkExpiresDuration: c.EmailVerifyLinkExpiresDuration,
	})

	eventProcessor, err := outbox.New(log, producer, usecases, &outboxConfig)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	grpcController := grpccontroller.New(log, usecases)

	kafkaConsumerController := kafkacontroller.New(log, &ku{})

	consumerGroup, err := kafka.NewConsumerGroup(&kafka.ConsumerGroupConfig{
		Brokers: c.KafkaBrokers,
		//TODO:change topics for real
		Topics:              []string{"permissions.events"},
		GroupID:             "users-service-app",
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
		grpcController,
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
