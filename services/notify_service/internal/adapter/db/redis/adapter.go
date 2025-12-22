package redisadapter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain"
	"github.com/redis/go-redis/v9"
)

var ErrCloseTimeout = errors.New("the time to safely terminate the connection to the redis has expired")

type Config struct {
	Addr     string
	Password string
	DB       int
}

type redisAdapter struct {
	rdb *redis.Client
	log *slog.Logger
}

func New(cfg *Config, log *slog.Logger) (*redisAdapter, error) {

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

	return &redisAdapter{
		rdb: rdb,
		log: log,
	}, nil

}

func (r *redisAdapter) SaveUserIDByChatID(ctx context.Context, chatID int64, userID string) error {

	const op = "adapter.redis.SaveUserIDByChatID"

	_, err := r.rdb.Set(ctx, chatUserKey(chatID), userID, 0).Result()

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil
}

func (r *redisAdapter) SaveChatIDByUserID(ctx context.Context, chatID int64, userID string) error {

	const op = "adapter.redis.SaveChatIDByUserID"

	_, err := r.rdb.Set(ctx, userChatKey(userID), chatID, 0).Result()

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil
}

func (r *redisAdapter) GetUserIDByChatID(ctx context.Context, chatID int64) (string, error) {

	const op = "adapter.redis.GetChatIDByUserID"

	userID, err := r.rdb.Get(ctx, chatUserKey(chatID)).Result()

	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("%s: %w: %v", op, adapter.ErrNotFound, err)
		}
		return "", fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return userID, nil
}

func (r *redisAdapter) GetChatIDByUserID(ctx context.Context, userID string) (int64, error) {

	const op = "adapter.redis.GetChatIDByUserID"

	chatIDstr, err := r.rdb.Get(ctx, userChatKey(userID)).Result()

	if err != nil {
		if err == redis.Nil {
			return 0, fmt.Errorf("%s: %w: %v", op, adapter.ErrNotFound, err)
		}
		return 0, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	chatID, err := strconv.ParseInt(chatIDstr, 10, 64)

	if err != nil {
		return 0, fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return chatID, nil
}

func (r *redisAdapter) SaveEventID(ctx context.Context, eventID string) error {
	const op = "adapter.redis.SaveEventID"

	_, err := r.rdb.Set(ctx, eventKey(eventID), domain.HandledEventStatus, 0).Result()

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return nil
}

func (r *redisAdapter) FindEvent(ctx context.Context, eventID string) (string, error) {

	const op = "adapter.redis.GetChatIDByUserID"

	status, err := r.rdb.Get(ctx, eventKey(eventID)).Result()

	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("%s: %w: %v", op, adapter.ErrNotFound, err)
		}
		return "", fmt.Errorf("%s: %w: %v", op, adapter.ErrInternal, err)
	}

	return status, nil
}

func (r *redisAdapter) Stop(ctx context.Context) error {
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
