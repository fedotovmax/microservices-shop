package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/pkg/utils/hashing"
)

func (u *usecases) CreateUser(ctx context.Context, meta *inputs.MetaParams, in *inputs.CreateUserInput) (string, error) {

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

		hashedPassword, err := hashing.Password(in.GetPassword())

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		in.SetPassword(hashedPassword)

		createUserResult, err = u.s.CreateUser(txCtx, in)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		link, err := u.s.CreateEmailVerifyLink(txCtx, createUserResult.ID)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		userCreatedPayload := events.UserCreatedEventPayload{
			EmailVerifyLinkValidityPeriod: link.ValidityPeriod,
			ID:                            createUserResult.ID,
			EmailVerifyLink:               link.Link,
			Email:                         createUserResult.Email,
			Locale:                        meta.GetLocale(),
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

		_, err = u.s.CreateEvent(txCtx, userCreatedEventIn)

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
