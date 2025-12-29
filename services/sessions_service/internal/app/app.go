package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/kafka-lib/kafka"
	"github.com/fedotovmax/microservices-shop-protos/events"
	jwtadapter "github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter/auth/jwt"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter/db/postgres"
	grpcadapter "github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter/grpc"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/config"
	grpccontroller "github.com/fedotovmax/microservices-shop/sessions_service/internal/controller/grpc_controller"
	kafkacontroller "github.com/fedotovmax/microservices-shop/sessions_service/internal/controller/kafka_controller"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/usecase"
	"github.com/fedotovmax/microservices-shop/sessions_service/pkg/logger"
	"github.com/fedotovmax/pgxtx"
)

type App struct {
	c             *config.AppConfig
	log           *slog.Logger
	postgres      postgres.PostgresPool
	grpc          *grpcadapter.Server
	consumerGroup kafka.ConsumerGroup
}

func New(c *config.AppConfig, log *slog.Logger) (*App, error) {

	const op = "app.New"

	l := log.With(slog.String("op", op))

	poolConnectCtx, poolConnectCtxCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer poolConnectCtxCancel()

	postgresPool, err := postgres.NewConnection(poolConnectCtx, &postgres.ConnectionConfig{
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

	postgresAdapter := postgres.NewAdapter(ex, log)

	jwtAdapter := jwtadapter.New(&jwtadapter.Config{
		AccessTokenExpDuration: c.AccessTokenExpDuration,
		AccessTokenSecret:      c.AccessTokenSecret,
	})

	usecases := usecase.New(log, txManager, jwtAdapter, postgresAdapter, c.RefreshTokenExpDuration)

	grpcController := grpccontroller.New(log, usecases)

	kafkaConsumerController := kafkacontroller.New(log, usecases)

	consumerGroup, err := kafka.NewConsumerGroup(&kafka.ConsumerGroupConfig{
		Brokers: c.KafkaBrokers,
		//TODO:change topics for real
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
		grpcController,
	)

	return &App{
			c:             c,
			log:           log,
			grpc:          grpcServer,
			postgres:      postgresPool,
			consumerGroup: consumerGroup,
		},
		nil
}

func (a *App) Run(cancel context.CancelFunc) {

	const op = "app.Run"

	log := a.log.With(slog.String("op", op))

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

	err := a.grpc.Stop(ctx)

	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("grpc server stopped successfully!")
	}

	err = a.consumerGroup.Stop(ctx)

	if err != nil {
		log.Error("error when stop consumer group", logger.Err(err))
	} else {
		log.Info("consumer group stopped successfully")
	}

	err = a.postgres.Stop(ctx)

	if err != nil {
		log.Error("error when stop postgres connection", logger.Err(err))
	} else {
		log.Info("postgres connection stopped successfully")
	}
}
