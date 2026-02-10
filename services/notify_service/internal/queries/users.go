package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/ports"
)

type Users interface {
	FindByChatID(ctx context.Context, chatID int64) (string, error)
}

type users struct {
	storage ports.UsersStorage
}

func NewUsers(s ports.UsersStorage) Users {
	return &users{
		storage: s,
	}
}

func (q *users) FindByChatID(ctx context.Context, chatID int64) (string, error) {

	userID, err := q.storage.FindByChatID(ctx, chatID)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return "", fmt.Errorf("%w: %v", errs.ErrUserIDNotFound, err)
		}
		return "", fmt.Errorf("%w", err)
	}

	return userID, nil
}
