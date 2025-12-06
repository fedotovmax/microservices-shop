package app

import (
	"context"
	"fmt"
	"log/slog"

	grpcadapter "github.com/fedotovmax/microservices-shop/api-gateway/internal/adapter/client/grpc"
	httpadapter "github.com/fedotovmax/microservices-shop/api-gateway/internal/adapter/http"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/controller"
	"github.com/fedotovmax/microservices-shop/api-gateway/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	HttpPort      uint16
	UsersGRPCAddr string
}

type service interface {
	Stop(ctx context.Context) error
}

type App struct {
	c           Config
	log         *slog.Logger
	http        *httpadapter.Server
	usersClient service
}

func New(log *slog.Logger, c Config) (*App, error) {

	const op = "app.New"

	r := chi.NewRouter()

	usersClient, err := grpcadapter.NewUsersClient(c.UsersGRPCAddr)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	customerController := controller.NewCustomerController(r, log, usersClient.RPC)

	customerController.Register()

	httpServer := httpadapter.NewHTTPAdapter(httpadapter.HTTPServerConfig{
		Port: c.HttpPort,
	}, r)

	return &App{
		c:           c,
		log:         log,
		http:        httpServer,
		usersClient: usersClient,
	}, nil
}

func (a *App) Run(cancel context.CancelFunc) {
	const op = "app.MustRun"

	log := a.log.With(slog.String("op", op))

	log.Info("Try to start HTTP server", slog.String("port", fmt.Sprintf("%d", a.c.HttpPort)))

	go func() {
		if err := a.http.Start(); err != nil {
			log.Error("Cannot start http server", logger.Err(err))
			log.Error("Signal to shutdown")
			cancel()
			return
		}
	}()
}

func (a *App) Stop(ctx context.Context) {

	const op = "app.Stop"

	log := a.log.With(slog.String("op", op))

	if err := a.http.Stop(ctx); err != nil {
		log.Error("Error when shutdown http server", logger.Err(err))
	} else {
		log.Info("HTTP server stopped successfully!")
	}

	//TODO: parallel close any grpc clients
	if err := a.usersClient.Stop(ctx); err != nil {
		log.Error("Error when stop GRPC users client", logger.Err(err))
	} else {
		log.Info("GRPC users client stopped successfully!")
	}
}
