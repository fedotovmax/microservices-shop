package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	_ "github.com/fedotovmax/microservices-shop/api-gateway/docs"
	grpcadapter "github.com/fedotovmax/microservices-shop/api-gateway/internal/adapter/client/grpc"
	httpadapter "github.com/fedotovmax/microservices-shop/api-gateway/internal/adapter/http"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/config"
	customercontroller "github.com/fedotovmax/microservices-shop/api-gateway/internal/controller/customer_controller"
	"github.com/fedotovmax/microservices-shop/api-gateway/pkg/logger"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type App struct {
	c                 *config.AppConfig
	log               *slog.Logger
	http              *httpadapter.Server
	lifesycleServices []*service
}

func New(log *slog.Logger, c *config.AppConfig) (*App, error) {
	const op = "app.New"

	r := chi.NewRouter()

	r.Handle("/swagger/*", httpSwagger.WrapHandler)

	usersClient, err := grpcadapter.NewUsersClient(c.UsersClientAddr)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	sessionsClient, err := grpcadapter.NewSessionsClient(c.SessionsClientAddr)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	customerController := customercontroller.New(r, log, usersClient.RPC, sessionsClient.RPC)

	customerController.Register()

	httpServer := httpadapter.NewHTTPAdapter(httpadapter.HTTPServerConfig{
		Port: c.Port,
	}, r)

	lifesycleServices := []*service{
		newService("users grpc client", usersClient),
		newService("sessions grpc client", sessionsClient),
	}

	return &App{
		c:                 c,
		log:               log,
		http:              httpServer,
		lifesycleServices: lifesycleServices,
	}, nil
}

func (a *App) Run(cancel context.CancelFunc) {
	const op = "app.MustRun"

	log := a.log.With(slog.String("op", op))

	log.Info("Try to start HTTP server", slog.String("port", fmt.Sprintf("%d", a.c.Port)))

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

	stopErrorsChan := make(chan error, len(a.lifesycleServices))

	wg := &sync.WaitGroup{}

	for _, s := range a.lifesycleServices {
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

	stopErrors := make([]error, 0, len(a.lifesycleServices))

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
