package chat

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis"
	goredis "github.com/redis/go-redis/v9"
)

func (r *redisDb) FindByUID(ctx context.Context, uid string) (int64, error) {

	const op = "adapters.redis.chat.FindByUID"

	chatIDstr, err := r.rdb.Get(ctx, redis.UserChatKey(uid)).Result()

	if err != nil {
		if err == goredis.Nil {
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
