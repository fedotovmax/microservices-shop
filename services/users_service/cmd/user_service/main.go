package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/fedotovmax/microservices-shop/users_service/internal/app"
	"github.com/fedotovmax/microservices-shop/users_service/internal/config"
	"github.com/fedotovmax/microservices-shop/users_service/internal/infra/logger"
)

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

	application, err := app.New(app.Config{
		DBURL:        cfg.DBUrl,
		GRPCPort:     uint16(cfg.Port),
		KafkaBrokers: cfg.KafkaBrokers,
	}, log)

	if err != nil {
		panic(err)
	}

	sig, triggerSignal := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer triggerSignal()

	application.MustRun(triggerSignal)

	<-sig.Done()

	log.Info("Signal recieved, shutdown app")

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer shutdownCtxCancel()

	application.Stop(shutdownCtx)

}
