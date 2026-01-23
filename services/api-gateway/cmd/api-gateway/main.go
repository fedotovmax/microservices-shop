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
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/app"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/config"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"github.com/fedotovmax/microservices-shop/api-gateway/pkg/logger"
)

func setupLooger(env string) (*slog.Logger, error) {
	switch env {
	case keys.Development:
		return logger.NewDevelopmentHandler(slog.LevelDebug), nil
	case keys.Production:
		return logger.NewProductionHandler(slog.LevelWarn), nil
	default:
		return nil, envconfig.ErrInvalidAppEnv
	}
}

// @title Swagger Documentation for microservices shop API Gateway
// @version 1.0
// @description Swagger Documentation for API Gateway of microservices shop project
// @contact.name Fedotv Maxim (developer)
// @contact.email f3d0t0tvmax@yandex.ru
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	application, err := app.New(log, cfg)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	sig, triggerSignal := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer triggerSignal()

	translationsDir := path.Join(workdir, cfg.TranslationPath)

	err = i18n.Local.Load(translationsDir)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	application.Run(triggerSignal)

	<-sig.Done()

	log.Info("Signal recieved, shutdown app")

	shutdownContext, shutdownCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer shutdownCancel()

	application.Stop(shutdownContext)
}
