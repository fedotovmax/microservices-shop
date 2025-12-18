package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/fedotovmax/envconfig"
	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/app"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/config"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/notify_service/pkg/logger"
)

func setupLooger(env string) (*slog.Logger, error) {
	switch env {
	case keys.Development:
		return logger.NewDevelopmentHandler(slog.LevelInfo), nil
	case keys.Production:
		return logger.NewProductionHandler(slog.LevelWarn), nil
	default:
		return nil, envconfig.ErrInvalidAppEnv
	}
}

func main() {
	cfg, err := config.LoadAppConfig()

	if err != nil {
		logger.GetFallback().Error(err.Error())
		os.Exit(1)
	}

	log, err := setupLooger(cfg.Env)

	if err != nil {
		logger.GetFallback().Error(err.Error())
		os.Exit(1)
	}

	workdir, err := os.Getwd()

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	translationsDir := path.Join(workdir, cfg.TranslationPath)
	err = i18n.Local.Load(translationsDir)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	application, err := app.New(&app.Config{
		KafkaBrokers:  cfg.KafkaBrokers,
		TgBotToken:    cfg.TgBotToken,
		RedisAddr:     cfg.RedisAddr,
		RedisPassword: cfg.RedisPassword,
	}, log)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	sig, triggerSignal := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer triggerSignal()

	application.Run()

	<-sig.Done()

	log.Info("Signal recieved, shutdown app")

	//TODO: real time out for shutdown
	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer shutdownCtxCancel()

	application.Stop(shutdownCtx)
}
