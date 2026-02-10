package events

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain"
)

func (r *redisDb) Save(ctx context.Context, eventID string) error {
	const op = "adapters.redis.events.Save"

	_, err := r.rdb.Set(ctx, redis.EventKey(eventID), domain.HandledEventStatus, 0).Result()

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil
}
