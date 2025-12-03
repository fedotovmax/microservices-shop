package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/fedotovmax/microservices-shop/api-gateway/internal/app"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/config"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/infra/logger"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
)

func mustSetupLooger(env string) *slog.Logger {
	switch env {
	case keys.Development:
		return logger.NewDevelopmentHandler()
	case keys.Production:
		return logger.NewProductionHandler()
	default:
		panic("unsopported app env for logger")
	}
}

func main() {

	cfg := config.MustLoadAppConfig()

	log := mustSetupLooger(cfg.Env)

	application, err := app.New(log, app.Config{
		HttpPort:      cfg.Port,
		UsersGRPCAddr: cfg.UsersClientAddr,
	})

	if err != nil {
		panic(err)
	}

	sig, triggerSignal := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer triggerSignal()

	application.MustRun(triggerSignal)

	<-sig.Done()

	log.Info("Signal recieved, shutdown app")

	shutdownContext, shutdownCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer shutdownCancel()

	application.Stop(shutdownContext)
}
