package redisadapter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCloseTimeout = errors.New("the time to safely terminate the connection to the redis has expired")

type Config struct {
	Addr     string
	Password string
	DB       int
}

type redisDb struct {
	rdb *redis.Client
	log *slog.Logger
}

func New(cfg *Config, log *slog.Logger) (*redisDb, error) {

	const op = "adapter.redis.New"

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	pingCtx, cancelPngCtx := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelPngCtx()

	_, err := rdb.Ping(pingCtx).Result()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &redisDb{
		rdb: rdb,
		log: log,
	}, nil

}

func (r *redisDb) GetClient() *redis.Client {
	return r.rdb
}

func (r *redisDb) Stop(ctx context.Context) error {
	op := "adapter.redis.Stop"

	done := make(chan error, 1)

	go func() {
		err := r.rdb.Close()
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%s: %w: %v", op, ErrCloseTimeout, ctx.Err())
	}
}
