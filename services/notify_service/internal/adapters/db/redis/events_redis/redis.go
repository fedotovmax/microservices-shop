package eventsredis

import (
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type redisDb struct {
	log *slog.Logger
	rdb *redis.Client
}

func New(log *slog.Logger, rdb *redis.Client) *redisDb {
	return &redisDb{
		log: log,
		rdb: rdb,
	}
}
