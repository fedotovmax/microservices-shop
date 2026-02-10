package events

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis"
	goredis "github.com/redis/go-redis/v9"
)

func (r *redisDb) FindByID(ctx context.Context, id string) (string, error) {

	const op = "adapters.redis.events.FindByID"

	status, err := r.rdb.Get(ctx, redis.EventKey(id)).Result()

	if err != nil {
		if err == goredis.Nil {
			return "", fmt.Errorf("%s: %w: %v", op, adapters.ErrNotFound, err)
		}
		return "", fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return status, nil
}
