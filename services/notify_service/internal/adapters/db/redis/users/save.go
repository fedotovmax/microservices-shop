package users

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis"
)

func (r *redisDb) Save(ctx context.Context, chatID int64, userID string) error {

	const op = "adapters.redis.users.Save"

	_, err := r.rdb.Set(ctx, redis.ChatUserKey(chatID), userID, 0).Result()

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil
}
