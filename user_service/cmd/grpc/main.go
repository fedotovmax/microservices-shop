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

	cfg, err := config.New()

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	poolConnectCtx, poolConnectCtxCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer poolConnectCtxCancel()

	postgresPool, err := infraPostgres.New(poolConnectCtx, cfg.DBUrl)

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("Successfully created db pool and connected!")

	txManager, err := pgxtx.Init(postgresPool)

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	ex := txManager.GetExtractor()

	userPostgres := adapterPostgres.NewUserPostgres(ex)

	newUser := domain.CreateUser{
		Email:     "makc@mail.ru",
		FirstName: "Maxim",
		LastName:  "Ivanov",
	}

	createUserCtx, createUserCtxCancel := context.WithTimeout(context.Background(), time.Second)
	defer createUserCtxCancel()
	userId, err := userPostgres.CreateUser(createUserCtx, newUser)

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("New User:", slog.String("id", userId))

	tcplistener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))

	if err != nil {
		slog.Error("Error when create net.Listen:", slog.String("error", err.Error()))
		os.Exit(1)
	}

	server := grpc.NewServer()

	svc := &service{}

	usersvc.RegisterUserServiceServer(server, svc)

	sigCtx, sigCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer sigCancel()

	go func() {
		slog.Info("Starting grpc server on port:", slog.Int("port", cfg.Port))
		if err := server.Serve(tcplistener); err != nil {
			slog.Error("Error when server grpc server:", slog.String("error", err.Error()))
			sigCancel()
			return
		}
	}()

	<-sigCtx.Done()

	slog.Info("Signal recieved, shutdown app")

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*15)
	defer shutdownCtxCancel()

	server.GracefulStop()

	postgresPool.GracefulStop(shutdownCtx)

	slog.Info("All resources are closed, exit app")

}
