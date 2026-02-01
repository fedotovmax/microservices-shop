package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/kafka-lib/kafka"
	"github.com/fedotovmax/microservices-shop-protos/events"
	redisadapter "github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis"
	chatredis "github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis/chat_redis"
	eventsredis "github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis/events_redis"
	usersredis "github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis/users_redis"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/telegram"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/config"
	"github.com/go-telegram/bot"

	kafkacontroller "github.com/fedotovmax/microservices-shop/notify_service/internal/controller/kafka_controller"
	tgbotcontroller "github.com/fedotovmax/microservices-shop/notify_service/internal/controller/tgbot_controller"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/usecase"
	"github.com/fedotovmax/microservices-shop/notify_service/pkg/logger"
)

type TGBot interface {
	Start()
	Stop()
}

type RedisAdapter interface {
	Stop(ctx context.Context) error
}

type App struct {
	c             *config.AppConfig
	log           *slog.Logger
	tgBot         TGBot
	redisAdapter  RedisAdapter
	consumerGroup kafka.ConsumerGroup
}

func New(c *config.AppConfig, log *slog.Logger) (*App, error) {

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

	opts := []bot.Option{}

	tgbot, err := telegram.New(&telegram.Config{
		Token:   c.TgBotToken,
		Options: opts,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	redisClient := redisAdapter.GetClient()

	usersRedis := usersredis.New(log, redisClient)
	chatRedis := chatredis.New(log, redisClient)
	eventsRedis := eventsredis.New(log, redisClient)

	usecases := usecase.New(log, usersRedis, chatRedis, eventsRedis, tgbot)

	kafkaConsumerController := kafkacontroller.NewKafkaController(log, usecases, &kafkacontroller.Config{
		CustomerSiteURL:                c.CustomerSiteURL,
		CustomerSiteURLEmailVerifyPath: c.CustomerSiteURLEmailVerifyPath,
	})

	tgBotController := tgbotcontroller.NewTgBotController(log, usecases, tgbot)

	tgBotController.Register()

	consumerGroup, err := kafka.NewConsumerGroup(&kafka.ConsumerGroupConfig{
		Brokers:             c.KafkaBrokers,
		Topics:              []string{events.USER_EVENTS, events.SESSION_EVENTS},
		GroupID:             "notify-service-app",
		SleepAfterRebalance: time.Second * 2,
		AutoCommit:          true,
	}, log, kafkaConsumerController)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &App{
			c:             c,
			log:           log,
			consumerGroup: consumerGroup,
			tgBot:         tgbot,
			redisAdapter:  redisAdapter,
		},
		nil

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

	err = a.redisAdapter.Stop(ctx)

	if err != nil {
		log.Error("error when stop redis client", logger.Err(err))
	} else {
		log.Info("redis client stopped successfully")
	}

}
