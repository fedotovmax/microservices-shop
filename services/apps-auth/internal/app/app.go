package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/fedotovmax/microservices-shop/apps-auth/internal/adapter/db/migrations"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/adapter/db/redisadapter"
	grpcadapter "github.com/fedotovmax/microservices-shop/apps-auth/internal/adapter/grpc"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/config"
	grpccontroller "github.com/fedotovmax/microservices-shop/apps-auth/internal/controller/grpc_controller"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/interceptors"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/usecase"
	"github.com/fedotovmax/microservices-shop/apps-auth/pkg/logger"
	"google.golang.org/grpc"
)

type App struct {
	c     *config.AuthAppConfig
	log   *slog.Logger
	grpc  *grpcadapter.Server
	redis *redisadapter.Rdb
}

func New(c *config.AuthAppConfig, log *slog.Logger) (*App, error) {

	const op = "app.New"

	l := log.With(slog.String("op", op))

	redisAdapter, err := redisadapter.New(&redisadapter.Config{
		Addr:     c.RedisAddr,
		Password: c.RedisPassword,
	}, log)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	l.Info("redis client successfully connected")

	usecases := usecase.NewUsecases(log, redisAdapter, &usecase.Config{
		TokenSecret:      c.TokenSecret,
		TokenExpDuration: c.TokenExpDuration,
		Issuer:           c.Issuer,
	})

	grpcController := grpccontroller.New(log, usecases)

	grpcServer := grpcadapter.New(
		grpcadapter.Config{
			Addr: fmt.Sprintf(":%d", c.Port),
		},
		grpcController,
		grpc.UnaryInterceptor(interceptors.AdminSecretInterceptor(c.AdminSecret)),
	)

	return &App{
			c:     c,
			log:   log,
			grpc:  grpcServer,
			redis: redisAdapter,
		},
		nil
}

func (a *App) Run(cancel context.CancelFunc) {

	const op = "app.Run"

	log := a.log.With(slog.String("op", op))

	log.Info("Applying redis migrations")

	migratorctx, migratorctxcancel := context.WithTimeout(context.Background(), time.Second*10)
	defer migratorctxcancel()

	err := migrations.ApplyRedisMigrations(migratorctx, a.c.AdminSecret, a.redis)

	if err != nil {
		log.Error("Error when apply redis migrations", logger.Err(err))
		log.Error("Signal to shutdown")
		cancel()
		return
	}

	log.Info("Redis migrations have been successfully applied")

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
		newService("redis", a.redis),
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
