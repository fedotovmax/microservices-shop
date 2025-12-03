package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	adapterPostgres "github.com/fedotovmax/microservices-shop/users_service/internal/adapter/postgres"
	"github.com/fedotovmax/microservices-shop/users_service/internal/controller"
	infraPostgres "github.com/fedotovmax/microservices-shop/users_service/internal/infra/db/postgres"
	"github.com/fedotovmax/microservices-shop/users_service/internal/infra/logger"
	"github.com/fedotovmax/microservices-shop/users_service/internal/infra/queues/kafka"
	grpcserver "github.com/fedotovmax/microservices-shop/users_service/internal/infra/server/grpc"
	"github.com/fedotovmax/microservices-shop/users_service/internal/usecase"
	"github.com/fedotovmax/outbox"

	"github.com/fedotovmax/pgxtx"
)

type Config struct {
	DBURL        string
	GRPCPort     uint16
	KafkaBrokers []string
}

type App struct {
	c        Config
	log      *slog.Logger
	postgres infraPostgres.PostgresPool
	grpc     *grpcserver.Server
	event    *outbox.Outbox
	producer kafka.Producer
}

func New(c Config, log *slog.Logger) (*App, error) {

	const op = "app.New"

	l := log.With(slog.String("op", op))

	poolConnectCtx, poolConnectCtxCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer poolConnectCtxCancel()

	postgresPool, err := infraPostgres.New(poolConnectCtx, c.DBURL)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	l.Info("Successfully created db pool and connected!")

	txManager, err := pgxtx.Init(postgresPool, log.With(slog.String("op", "transaction.manager")))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ex := txManager.GetExtractor()

	postgresAdapter := adapterPostgres.NewPostgresAdapter(ex)

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

	eventProcessor, err := outbox.New(log, producer, txManager, ex, outboxConfig)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	eventSender := eventProcessor.GetEventSender()

	useceses := usecase.NewUsecases(postgresAdapter, txManager, eventSender)
	grpcController := controller.NewGRPCController(log, useceses)

	grpcServer := grpcserver.NewGRPCServer(
		grpcserver.Config{
			Addr: fmt.Sprintf(":%d", c.GRPCPort),
		},
		grpcController,
	)

	return &App{
			c:        c,
			log:      log,
			grpc:     grpcServer,
			postgres: postgresPool,
			event:    eventProcessor,
			producer: producer},
		nil
}

func (a *App) MustRun(cancel ...context.CancelFunc) {

	const op = "app.MustRun"

	log := a.log.With(slog.String("op", op))

	//a.event.Start()

	log.Info("event processor starting")

	log.Info("Try to start GRPC server:", slog.String("port", fmt.Sprintf("%d", a.c.GRPCPort)))

	go func() {
		if err := a.grpc.Start(); err != nil {
			if len(cancel) > 0 {
				log.Error("Cannot start grpc server", logger.Err(err))
				log.Error("Signal to shutdown")
				cancel[0]()
				return
			}
			panic(fmt.Errorf("%s: %w", op, err))
		}
	}()

}

func (a *App) Stop(ctx context.Context) {

	const op = "app.Stop"

	log := a.log.With(slog.String("op", op))

	stopErrosChan := make(chan error)

	err := a.grpc.Stop(ctx)

	go func() {
		if err != nil {
			stopErrosChan <- err
		} else {
			log.Info("grpc server stopped successfully!")
		}
	}()

	lifesycleServices := []*service{
		//newService("event processor", a.event),
		newService("kafka producer", a.producer),
		newService("postgres", a.postgres),
	}

	wg := &sync.WaitGroup{}

	for _, s := range lifesycleServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := s.Stop(ctx)
			if err != nil {
				stopErrosChan <- err
			} else {
				log.Info("successfully stopped:", slog.String("service", s.Name()))
			}
		}()
	}

	go func() {
		wg.Wait()
		close(stopErrosChan)
	}()

	stopErrors := make([]error, 0)

	for err := range stopErrosChan {
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
