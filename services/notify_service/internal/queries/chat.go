package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/ports"
)

type Chat interface {
	FindByUserID(ctx context.Context, uid string) (int64, error)
}

type chat struct {
	storage ports.ChatStorage
}

func NewChat(s ports.ChatStorage) Chat {
	return &chat{
		storage: s,
	}
}

func (q *chat) FindByUserID(ctx context.Context, uid string) (int64, error) {

	chatID, err := q.storage.FindByUID(ctx, uid)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return 0, fmt.Errorf("%w: %v", errs.ErrChatIDNotFound, err)
		}
		return 0, fmt.Errorf("%w", err)
	}

	return chatID, nil
}
