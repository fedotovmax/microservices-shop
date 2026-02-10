package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/ports"
)

type Events interface {
	FindByID(ctx context.Context, id string) (string, error)
}

type events struct {
	storage ports.EventsStorage
}

func NewEvents(s ports.EventsStorage) Events {
	return &events{
		storage: s,
	}
}

func (q *events) FindByID(ctx context.Context, id string) (string, error) {

	chatID, err := q.storage.FindByID(ctx, id)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return "", fmt.Errorf("%w: %v", errs.ErrEventNotFound, err)
		}
		return "", fmt.Errorf("%w", err)
	}

	return chatID, nil
}
