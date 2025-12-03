package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/utils/hashing"
	"github.com/fedotovmax/outbox"
)

func (u *usecases) CreateUser(ctx context.Context, d *domain.CreateUserInput) (string, error) {
	const op = "usecase.user.CreateUser"

	var userId string

	hashedPassword, err := hashing.Password(d.GetPassword())

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	d.SetPasswordHash(hashedPassword)

	err = u.txm.Wrap(ctx, func(txCtx context.Context) error {
		var err error
		userId, err = u.s.Create(txCtx, d)

		if err != nil {
			return err
		}

		payload := domain.UserCreatedEvent{ID: userId}

		b, err := json.Marshal(payload)

		if err != nil {
			return err
		}

		u.es.AddNewEvent(txCtx, outbox.CreateEvent{
			AggregateID: userId,
			Topic:       events.USER_EVENTS,
			Type:        "user.created",
			Payload:     b,
		})

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil

}
