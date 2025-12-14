package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/pkg/utils/hashing"
)

func (u *usecases) CreateUser(ctx context.Context, in *inputs.CreateUserInput) (string, error) {
	const op = "usecase.CreateUser"

	var userId string

	err := u.txm.Wrap(ctx, func(txCtx context.Context) error {

		var err error

		_, err = u.FindUserByEmail(txCtx, in.GetEmail())

		if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
			return fmt.Errorf("%s: %w", op, err)
		}

		if err == nil {
			return fmt.Errorf("%s: %w", op, errs.ErrUserAlreadyExists)
		}

		hashedPassword, err := hashing.Password(in.GetPassword())

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		in.SetPassword(hashedPassword)

		userId, err = u.s.CreateUser(txCtx, in)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		payload := events.UserCreatedEventPayload{ID: userId}

		b, err := json.Marshal(payload)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		eventIn := inputs.NewCreateEventInput()

		eventIn.SetAggregateID(userId)
		eventIn.SetTopic(events.USER_EVENTS)
		eventIn.SetType(events.USER_CREATED)
		eventIn.SetPayload(b)

		_, err = u.s.CreateEvent(txCtx, eventIn)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return userId, nil

}
