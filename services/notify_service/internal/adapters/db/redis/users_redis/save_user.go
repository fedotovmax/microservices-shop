package usersredis

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters"
	redisadapter "github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis"
)

func (r *redisDb) SaveUserIDByChatID(ctx context.Context, chatID int64, userID string) error {

	const op = "adapter.redis.users.SaveUserIDByChatID"

	_, err := r.rdb.Set(ctx, redisadapter.ChatUserKey(chatID), userID, 0).Result()

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil
}
