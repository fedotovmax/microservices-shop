package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/fedotovmax/microservices-shop/apps-auth/internal/adapter/db/redisadapter"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/config"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/domain"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/utils"
	"github.com/fedotovmax/microservices-shop/apps-auth/pkg/logger"
)

func main() {
	cfg, err := config.LoadCliAuthAppConfig()

	if err != nil {
		logger.GetFallback().Error(err.Error())
		os.Exit(1)
	}

	log := logger.NewProductionHandler(slog.LevelInfo)

	rdb, err := redisadapter.New(&redisadapter.Config{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	}, log)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	newAppType := domain.ApplicationType(cfg.NewAppType)

	err = newAppType.IsValid()

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	app := domain.App{CreatedAt: time.Now().UTC(), Name: cfg.NewAppName, Type: newAppType}

	saveCtx, cancelSaveCtx := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelSaveCtx()

	secret, err := utils.CreateSecret()

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	secretHash := utils.CreateHash(secret)

	err = rdb.SaveApp(saveCtx, secretHash, &app)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	log.Info("New application created!", slog.String("secret", secret))

	rdbStopCtx, rdbStopCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer rdbStopCancel()

	rdb.Stop(rdbStopCtx)

	log.Info("CLI apps-auth app stopped")

}
