package usersredis

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters"
	redisadapter "github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis"
	"github.com/redis/go-redis/v9"
)

func (r *redisDb) GetUserIDByChatID(ctx context.Context, chatID int64) (string, error) {

	const op = "adapter.redis.GetChatIDByUserID"

	userID, err := r.rdb.Get(ctx, redisadapter.ChatUserKey(chatID)).Result()

	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("%s: %w: %v", op, adapters.ErrNotFound, err)
		}
		return "", fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return userID, nil
}
