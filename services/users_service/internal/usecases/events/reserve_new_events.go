package events

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/goutils/sliceutils"
	"github.com/fedotovmax/kafka-lib/outbox"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
)

func (u *usecases) ReserveNewEvents(ctx context.Context, limit int, reserveDuration time.Duration) ([]outbox.Event, error) {

	const op = "usecase.events.ReserveNewEvents"

	var events []*domain.Event

	err := u.txm.Wrap(ctx, func(txCtx context.Context) error {
		var err error
		events, err = u.storage.FindNewAndNotReservedEvents(txCtx, limit)

		if err != nil {
			return err
		}

		eventsIds := make([]string, len(events))

		for i := 0; i < len(events); i++ {
			eventsIds[i] = events[i].ID
		}

		err = u.storage.SetEventsReservedToByIDs(txCtx, eventsIds, reserveDuration)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sliceutils.SliceToSliceInterface[*domain.Event, outbox.Event](events), nil
}
