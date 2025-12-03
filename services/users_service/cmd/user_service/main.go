package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop/users_service/internal/app"
	"github.com/fedotovmax/microservices-shop/users_service/internal/config"
	"github.com/fedotovmax/microservices-shop/users_service/internal/infra/logger"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
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

	workdir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	cfg := config.MustLoadAppConfig()

	log := mustSetupLooger(cfg.Env)

	application, err := app.New(app.Config{
		DBURL:        cfg.DBUrl,
		GRPCPort:     cfg.Port,
		KafkaBrokers: cfg.KafkaBrokers,
	}, log)

	if err != nil {
		panic(err)
	}

	sig, triggerSignal := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer triggerSignal()

	translationsDir := path.Join(workdir, cfg.TranslationPath)

	err = i18n.Manager.Load(log, translationsDir)

	if err != nil {
		panic(err)
	}

	application.MustRun(triggerSignal)

	<-sig.Done()

	log.Info("Signal recieved, shutdown app")

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer shutdownCtxCancel()

	application.Stop(shutdownCtx)

}
