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
	adapterKafka "github.com/fedotovmax/microservices-shop/user_service/internal/adapter/kafka"
	adapterPostgres "github.com/fedotovmax/microservices-shop/user_service/internal/adapter/postgres"
	eventprocessor "github.com/fedotovmax/microservices-shop/user_service/internal/app/event-processor"
	"github.com/fedotovmax/microservices-shop/user_service/internal/config"
	"github.com/fedotovmax/microservices-shop/user_service/internal/domain"
	infraPostgres "github.com/fedotovmax/microservices-shop/user_service/internal/infra/db/postgres"
	"github.com/fedotovmax/microservices-shop/user_service/internal/infra/logger"
	infraKafka "github.com/fedotovmax/microservices-shop/user_service/internal/infra/queues/kafka"
	"github.com/fedotovmax/microservices-shop/user_service/internal/usecase"
	"github.com/fedotovmax/pgxtx"
	"google.golang.org/grpc"
)

type service struct {
	usersvc.UnimplementedUserServiceServer
}

func (s *service) CreateUser(ctx context.Context, req *usersvc.CreateUserRequest) (*usersvc.CreateUserResponse, error) {
	return nil, fmt.Errorf("some error")
}

func main() {
	cfg := config.MustLoadAppConfig()

	log := logger.MustNewLogger(cfg.Env)

	poolConnectCtx, poolConnectCtxCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer poolConnectCtxCancel()

	postgresPool, err := infraPostgres.New(poolConnectCtx, cfg.DBUrl)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	log.Info("Successfully created db pool and connected!")

	txManager, err := pgxtx.Init(postgresPool)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	ex := txManager.GetExtractor()

	produceInsatnce, err := infraKafka.NewAsyncProducer(cfg.KafkaBrokers)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	producerKafka := adapterKafka.NewProduceAdapter(produceInsatnce)

	// postgres adapters
	userPostgres := adapterPostgres.NewUserPostgres(ex)
	eventPostgres := adapterPostgres.NewEventPostgres(ex)

	// usecases
	userUsecase := usecase.NewUserUsecase(userPostgres, eventPostgres, txManager)
	eventUsecase := usecase.NewEventUsecase(eventPostgres, txManager)
	// TODO: get all params from config!
	eventProcessor := eventprocessor.New(log, producerKafka, eventUsecase, eventprocessor.Config{
		Limit:   50,
		Workers: 5,
	})

	createUserCtx, cancelCreateUserCtx := context.WithTimeout(context.Background(), time.Second)
	defer cancelCreateUserCtx()

	userId, err := userUsecase.CreateUser(createUserCtx, domain.CreateUser{Email: "makc-ivanov@mail.ru", FirstName: "Maxim", LastName: "Ivanov"})

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

	produceInsatnce.Close(shutdownCtx)

	log.Info("All resources are closed, exit app")

}
