package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain/errs"
)

func (u *usecases) IsNewEvent(ctx context.Context, eventID string) error {

	const op = "usecases.IsNewEvent"

	_, err := u.storage.FindEvent(ctx, eventID)

	if err != nil {
		if errors.Is(err, adapter.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return fmt.Errorf("%s: %w", op, errs.ErrEventAlreadyHandled)
}
