package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/usersvc"
	adapterPostgres "github.com/fedotovmax/microservices-shop/user_service/internal/adapter/postgres"
	"github.com/fedotovmax/microservices-shop/user_service/internal/config"
	"github.com/fedotovmax/microservices-shop/user_service/internal/domain"
	infraPostgres "github.com/fedotovmax/microservices-shop/user_service/internal/infra/db/postgres"
	"github.com/fedotovmax/microservices-shop/user_service/internal/infra/logger"
	"github.com/fedotovmax/microservices-shop/user_service/internal/infra/queues/kafka"
	"github.com/fedotovmax/microservices-shop/user_service/internal/usecase"
	"github.com/fedotovmax/outbox"
	"github.com/fedotovmax/pgxtx"
	"google.golang.org/grpc"
)

type service struct {
	usersvc.UnimplementedUserServiceServer
}

func mustSetupLooger(env string) *slog.Logger {
	switch env {
	case config.Development:
		return logger.NewDevelopmentHandler()
	case config.Production:
		return logger.NewProductionHandler()
	default:
		panic("unsopported app env for logger")
	}
}

func main() {
	cfg := config.MustLoadAppConfig()

	log := mustSetupLooger(cfg.Env)

	poolConnectCtx, poolConnectCtxCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer poolConnectCtxCancel()

	postgresPool, err := infraPostgres.New(poolConnectCtx, cfg.DBUrl)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	log.Info("Successfully created db pool and connected!")

	txManager, err := pgxtx.Init(postgresPool, log.With(slog.String("op", "transaction.manager")))

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	ex := txManager.GetExtractor()

	producer, err := kafka.NewAsyncProducer(cfg.KafkaBrokers)

	eventProcessor := outbox.New(log, producer, txManager, ex, outbox.Config{
		Limit:   50,
		Workers: 5,
	})

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	// postgres adapters
	userPostgres := adapterPostgres.NewUserPostgres(ex)

	// usecases
	userUsecase := usecase.NewUserUsecase(userPostgres, txManager, eventProcessor)
	// TODO: get all params from config!

	createUserCtx, cancelCreateUserCtx := context.WithTimeout(context.Background(), time.Second)
	defer cancelCreateUserCtx()

	userId, err := userUsecase.CreateUser(createUserCtx, domain.CreateUser{Email: "makc-ivanov@mail.ru"})

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	log.Info("User Created:", slog.String("user_id", userId))

	tcplistener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))

	if err != nil {
		log.Error("Error when create net.Listen:", logger.Err(err))
		os.Exit(1)
	}

	server := grpc.NewServer()

	svc := &service{}

	usersvc.RegisterUserServiceServer(server, svc)

	sigCtx, sigCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer sigCancel()

	eventProcessor.Start()
	log.Debug("eventProcessor starting")

	go func() {
		log.Info("Starting grpc server on port:", slog.Int("port", cfg.Port))
		if err := server.Serve(tcplistener); err != nil {
			log.Error("Error when server grpc server:", logger.Err(err))
			sigCancel()
			return
		}
	}()

	<-sigCtx.Done()

	log.Info("Signal recieved, shutdown app")

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*15)
	defer shutdownCtxCancel()

	server.GracefulStop()
	eventProcessor.Stop(shutdownCtx)

	postgresPool.GracefulStop(shutdownCtx)

	producer.Close(shutdownCtx)

	log.Info("All resources are closed, exit app")

}
