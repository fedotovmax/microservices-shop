package chatredis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters"
	redisadapter "github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis"
	"github.com/redis/go-redis/v9"
)

func (r *redisDb) GetChatIDByUserID(ctx context.Context, userID string) (int64, error) {

	const op = "adapter.redis.chat.GetChatIDByUserID"

	chatIDstr, err := r.rdb.Get(ctx, redisadapter.UserChatKey(userID)).Result()

	if err != nil {
		if err == redis.Nil {
			return 0, fmt.Errorf("%s: %w: %v", op, adapters.ErrNotFound, err)
		}
		return 0, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	chatID, err := strconv.ParseInt(chatIDstr, 10, 64)

	if err != nil {
		return 0, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return chatID, nil
}
