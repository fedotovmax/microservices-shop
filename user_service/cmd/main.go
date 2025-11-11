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
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type service struct {
	usersvc.UnimplementedUserServiceServer
}

func (s *service) CreateUser(ctx context.Context, req *usersvc.CreateUserRequest) (*usersvc.CreateUserResponse, error) {
	return nil, fmt.Errorf("some error")
}

func main() {

	poolConnectCtx, poolConnectCtxCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer poolConnectCtxCancel()
	postgresUrl := os.Getenv("DB_URL")
	if postgresUrl == "" {
		slog.Error("postgres db url not provided!")
		os.Exit(1)
	}

	pool, err := pgxpool.New(poolConnectCtx, postgresUrl)
	if err != nil {
		slog.Error("Cannot connect to postges:", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = pool.Ping(context.Background())

	if err != nil {
		slog.Error("Error when ping postgres:", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("Successfully connected to postgtes!")

	port := os.Getenv("PORT")

	if port == "" {
		slog.Error("port not provided!")
		os.Exit(1)
	}

	tcplistener, err := net.Listen("tcp", ":"+port)

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
		slog.Info("Starting grpc server on port:", slog.String("port", port))
		if err := server.Serve(tcplistener); err != nil {
			slog.Error("Error when server grpc server:", slog.String("error", err.Error()))
			sigCancel()
			return
		}
	}()

	<-sigCtx.Done()

	slog.Info("Signal recieved, shutdown app")

	server.GracefulStop()

	pool.Close()

	slog.Info("All resources are closed, exit app")

}
