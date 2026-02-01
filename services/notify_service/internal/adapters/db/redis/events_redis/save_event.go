package eventsredis

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters"
	redisadapter "github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain"
)

func (r *redisDb) SaveEventID(ctx context.Context, eventID string) error {
	const op = "adapter.redis.events.SaveEventID"

	_, err := r.rdb.Set(ctx, redisadapter.EventKey(eventID), domain.HandledEventStatus, 0).Result()

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil
}
