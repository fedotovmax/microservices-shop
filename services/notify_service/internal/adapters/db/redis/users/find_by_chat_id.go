package users

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis"
	goredis "github.com/redis/go-redis/v9"
)

func (r *redisDb) FindByChatID(ctx context.Context, chatID int64) (string, error) {

	const op = "adapters.redis.users.FindByChatID"

	userID, err := r.rdb.Get(ctx, redis.ChatUserKey(chatID)).Result()

	if err != nil {
		if err == goredis.Nil {
			return "", fmt.Errorf("%s: %w: %v", op, adapters.ErrNotFound, err)
		}
		return "", fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return userID, nil
}
