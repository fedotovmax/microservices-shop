package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/user_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/user_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/user_service/internal/ports"
	"github.com/fedotovmax/pgxtx"
)

type EventUsesace struct {
	ea  ports.EventAdapter
	txm pgxtx.Manager
}

func NewEventUsecase(ea ports.EventAdapter, txm pgxtx.Manager) *EventUsesace {

	return &EventUsesace{
		ea:  ea,
		txm: txm,
	}
}

func (e *EventUsesace) ConfirmFailed(ctx context.Context, ev *domain.FailedEvent) error {
	const op = "usecase.event.ConfirmFailed"

	err := e.txm.Wrap(ctx, func(txCtx context.Context) error {

		err := e.ea.RemoveReserve(txCtx, ev.ID)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (e *EventUsesace) ConfirmEvent(ctx context.Context, ev *domain.SuccessEvent) error {

	const op = "usecase.event.ConfirmEvent"

	err := e.txm.Wrap(ctx, func(txCtx context.Context) error {

		err := e.ea.RemoveReserve(txCtx, ev.ID)

		if err != nil {
			return err
		}

		err = e.ea.ChangeStatus(txCtx, ev.ID)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (e *EventUsesace) ReserveNewEvents(ctx context.Context, limit int, reserveDuration time.Duration) ([]*domain.Event, error) {

	const op = "usecase.event.ReserveNewEvents"

	var events []*domain.Event

	err := e.txm.Wrap(ctx, func(txCtx context.Context) error {
		var err error
		events, err = e.ea.FindNewAndNotReserved(txCtx, limit)

		if err != nil {
			return err
		}

		eventsIds := make([]string, len(events))

		for i := 0; i < len(events); i++ {
			eventsIds[i] = events[i].ID
		}

		err = e.ea.SetReservedToByIDs(txCtx, eventsIds, reserveDuration)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w: %v", op, errs.ErrNoNewEvents, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}
