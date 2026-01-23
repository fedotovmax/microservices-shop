package events

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
)

func (u *usecases) CreateEvent(ctx context.Context, d *inputs.CreateEvent) (string, error) {
	const op = "usecase.events.ConfirmFailedEvent"

	res, err := u.storage.CreateEvent(ctx, d)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}
