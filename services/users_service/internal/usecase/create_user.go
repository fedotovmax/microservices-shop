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

func (u *usecases) CreateUser(ctx context.Context, in *inputs.CreateUserInput, emailIn *inputs.EmailVerifyNotificationInput) (string, error) {
	const op = "usecase.CreateUser"

	var userID string

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

		userID, err = u.s.CreateUser(txCtx, in)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		link, err := u.s.CreateEmailVerifyLink(txCtx, userID)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		userCreatedPayload := events.UserCreatedEventPayload{ID: userID}

		userCreatedPayloadBytes, err := json.Marshal(userCreatedPayload)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		userCreatedEventIn := inputs.NewCreateEventInput()

		userCreatedEventIn.SetAggregateID(userID)
		userCreatedEventIn.SetTopic(events.USER_EVENTS)
		userCreatedEventIn.SetType(events.USER_CREATED)
		userCreatedEventIn.SetPayload(userCreatedPayloadBytes)

		_, err = u.s.CreateEvent(txCtx, userCreatedEventIn)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		emailVerifyNotificationPayload := events.EmailVerifyNotificationPayload{
			Email:       in.GetEmail(),
			Title:       emailIn.GetTitle(),
			Description: emailIn.GetDescription(),
			Link:        link.Link,
			Locale:      emailIn.GetLocale(),
		}

		emailVerifyNotificationPayloadBytes, err := json.Marshal(emailVerifyNotificationPayload)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		emailVerifyNotificationIn := inputs.NewCreateEventInput()
		emailVerifyNotificationIn.SetAggregateID(userID)
		emailVerifyNotificationIn.SetTopic(events.NOTIFICATIONS_EVENTS)
		emailVerifyNotificationIn.SetType(events.NOTIFICATIONS_EMAIL)
		emailVerifyNotificationIn.SetPayload(emailVerifyNotificationPayloadBytes)

		_, err = u.s.CreateEvent(txCtx, emailVerifyNotificationIn)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return userID, nil

}
