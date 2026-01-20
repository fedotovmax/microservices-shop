package redisadapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop/apps-auth/internal/adapter"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/domain"
	"github.com/redis/go-redis/v9"
)

var ErrCloseTimeout = errors.New("the time to safely terminate the connection to the redis has expired")

type Config struct {
	Addr     string
	Password string
	DB       int
}

type Rdb struct {
	rdb *redis.Client
	log *slog.Logger
}

func New(cfg *Config, log *slog.Logger) (*Rdb, error) {

	const op = "adapter.db.redis.New"

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

	return &Rdb{
		rdb: rdb,
		log: log,
	}, nil

}

func (r *Rdb) SaveApp(ctx context.Context, secretHash string, app *domain.App) error {

	const op = "adapter.db.redis.SaveApp"

	data, err := json.Marshal(app)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	err = r.rdb.Set(ctx, appKey(secretHash), data, 0).Err()

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil
}

func (r *Rdb) FindApp(ctx context.Context, secretHash string) (*domain.App, error) {
	const op = "adapter.db.redis.FindApp"

	appBytes, err := r.rdb.Get(ctx, appKey(secretHash)).Bytes()

	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	var app domain.App

	err = json.Unmarshal(appBytes, &app)

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return &app, nil
}

func (r *Rdb) DeleteApp(ctx context.Context, secretHash string) error {
	op := "adapter.db.redis.DeleteApp"

	err := r.rdb.Del(ctx, secretHash).Err()

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil

}

func (r *Rdb) Stop(ctx context.Context) error {
	op := "adapter.db.redis.Stop"

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
