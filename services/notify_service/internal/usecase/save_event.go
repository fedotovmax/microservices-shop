package usecase

import (
	"context"
	"fmt"
)

func (u *usecases) SaveEvent(ctx context.Context, eventID string) error {

	const op = "usecases.SaveEvent"

	err := u.eventsStorage.SaveEventID(ctx, eventID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
