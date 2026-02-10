package app

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/kafka-lib/kafka"
	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis"
	redisChat "github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis/chat"
	redisEvents "github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis/events"
	redisUsers "github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis/users"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/telegram"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/config"
	kafkaController "github.com/fedotovmax/microservices-shop/notify_service/internal/controller/kafka"
	telegramController "github.com/fedotovmax/microservices-shop/notify_service/internal/controller/telegram"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/queries"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/usecases"
	"github.com/go-telegram/bot"
)

func New(c *config.AppConfig, log *slog.Logger) (*App, error) {

	const op = "app.New"

	l := log.With(slog.String("op", op))

	rdb, err := redis.New(&redis.Config{
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

	redisClient := rdb.GetClient()

	usersRedis := redisUsers.New(log, redisClient)
	chatRedis := redisChat.New(log, redisClient)
	eventsRedis := redisEvents.New(log, redisClient)

	chatQuery := queries.NewChat(chatRedis)
	//usersQuery := queries.NewUsers(usersRedis)
	eventsQuery := queries.NewEvents(eventsRedis)

	isNewEvent := usecases.NewIsNewEventUsecase(log, eventsQuery)
	saveEvent := usecases.NewSaveEventUsecase(log, eventsRedis)
	saveChatUserPair := usecases.NewSaveChatUserPairUsecase(log, chatRedis, usersRedis)
	sendTgMessage := usecases.NewSendTgMessageUsecase(l, chatQuery, tgbot)

	kafkaConsumerController := kafkaController.New(
		log,
		sendTgMessage,
		saveEvent,
		isNewEvent,
		&kafkaController.Config{
			CustomerSiteURL:                c.CustomerSiteURL,
			CustomerSiteURLEmailVerifyPath: c.CustomerSiteURLEmailVerifyPath,
		},
	)

	tgBotController := telegramController.New(
		log,
		saveChatUserPair,
		tgbot,
	)

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
			redis:         rdb,
		},
		nil

}
