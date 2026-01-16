package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
)

func (u *usecases) CreateUser(ctx context.Context, in *inputs.CreateUserInput, locale string) (string, error) {

	const op = "usecase.CreateUser"

	var createUserResult *domain.UserPrimaryFields

	err := u.txm.Wrap(ctx, func(txCtx context.Context) error {

		var err error

		_, err = u.FindUserByEmail(txCtx, in.GetEmail())

		if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
			return fmt.Errorf("%s: %w", op, err)
		}

		if err == nil {
			return fmt.Errorf("%s: %w", op, errs.ErrUserAlreadyExists)
		}

		hashedPassword, err := hashPassword(in.GetPassword())

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		in.SetPassword(hashedPassword)

		createUserResult, err = u.s.users.CreateUser(txCtx, in)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		expiresAt := time.Now().Add(u.cfg.EmailVerifyLinkExpiresDuration).UTC()

		link, err := u.s.users.CreateEmailVerifyLink(txCtx, createUserResult.ID, expiresAt)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		userCreatedPayload := events.UserCreatedEventPayload{
			ID:                       createUserResult.ID,
			EmailVerifyLink:          link.Link,
			EmailVerifyLinkExpiresAt: link.LinkExpiresAt,
			Email:                    createUserResult.Email,
			Locale:                   locale,
		}

		userCreatedPayloadBytes, err := json.Marshal(userCreatedPayload)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		userCreatedEventIn := inputs.NewCreateEventInput()
		userCreatedEventIn.SetAggregateID(createUserResult.ID)
		userCreatedEventIn.SetTopic(events.USER_EVENTS)
		userCreatedEventIn.SetType(events.USER_CREATED)
		userCreatedEventIn.SetPayload(userCreatedPayloadBytes)

		_, err = u.s.events.CreateEvent(txCtx, userCreatedEventIn)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return createUserResult.ID, nil

}
