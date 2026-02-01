package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
)

func (u *usecases) CreateUser(ctx context.Context, in *inputs.CreateUserInput, locale string) (string, error) {

	const op = "usecase.users.CreateUser"

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

		createUserResult, err = u.usersStorage.CreateUser(txCtx, in)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		expiresAt := time.Now().Add(u.cfg.EmailVerifyLinkExpiresDuration).UTC()

		link, err := u.emailVerifyStorage.CreateEmailVerifyLink(txCtx, createUserResult.ID, expiresAt)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.SendUserCreatedEvent(txCtx, &sendUserCreatedEventParams{
			ID:     createUserResult.ID,
			Email:  createUserResult.Email,
			Locale: locale,
		})

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.SendEmalVerifyLinkAddedEvent(txCtx, &sendEmalVerifyLinkAddedEventParams{
			ID:            createUserResult.ID,
			Email:         createUserResult.Email,
			Link:          link.Link,
			LinkExpiresAt: link.LinkExpiresAt,
			Locale:        locale,
		})

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
