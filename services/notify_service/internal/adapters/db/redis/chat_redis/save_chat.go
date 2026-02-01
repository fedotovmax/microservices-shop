package chatredis

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters"
	redisadapter "github.com/fedotovmax/microservices-shop/notify_service/internal/adapters/db/redis"
)

func (r *redisDb) SaveChatIDByUserID(ctx context.Context, chatID int64, userID string) error {

	const op = "adapter.redis.chat.SaveChatIDByUserID"

	_, err := r.rdb.Set(ctx, redisadapter.UserChatKey(userID), chatID, 0).Result()

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil
}
