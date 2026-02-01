package eventsredis

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters"
	redisadapter "github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis"
	"github.com/redis/go-redis/v9"
)

func (r *redisDb) FindEvent(ctx context.Context, eventID string) (string, error) {

	const op = "adapter.redis.events.FindEvent"

	status, err := r.rdb.Get(ctx, redisadapter.EventKey(eventID)).Result()

	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("%s: %w: %v", op, adapters.ErrNotFound, err)
		}
		return "", fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return status, nil
}
