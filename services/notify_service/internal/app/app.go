package app

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/kafka-lib/kafka"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/config"
	"github.com/fedotovmax/microservices-shop/notify_service/pkg/logger"
)

type TGBot interface {
	Start()
	Stop()
}

type Service interface {
	Stop(ctx context.Context) error
}

type App struct {
	c             *config.AppConfig
	log           *slog.Logger
	tgBot         TGBot
	redis         Service
	consumerGroup kafka.ConsumerGroup
}

func (a *App) Run() {
	const op = "app.Run"

	log := a.log.With(slog.String("op", op))

	a.consumerGroup.Start()
	log.Info("consumer group starting")

	a.tgBot.Start()
	log.Info("trying to start telegram bot")

}

func (a *App) Stop(ctx context.Context) {

	const op = "app.Stop"

	log := a.log.With(slog.String("op", op))

	a.tgBot.Stop()
	log.Info("telegram bot stoppped")

	err := a.consumerGroup.Stop(ctx)

	if err != nil {
		log.Error("error when stop consumer group", logger.Err(err))
	} else {
		log.Info("consumer group stopped successfully")
	}

	err = a.redis.Stop(ctx)

	if err != nil {
		log.Error("error when stop redis client", logger.Err(err))
	} else {
		log.Info("redis client stopped successfully")
	}

}
