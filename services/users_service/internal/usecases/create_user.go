package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/ports"
	"github.com/fedotovmax/microservices-shop/users_service/internal/publisher"
	"github.com/fedotovmax/microservices-shop/users_service/internal/queries"
	"github.com/fedotovmax/microservices-shop/users_service/internal/utils"
	"github.com/fedotovmax/pgxtx"
)

type CreateUserUsecase struct {
	txm               pgxtx.Manager
	log               *slog.Logger
	cfg               *EmailConfig
	usersStorage      ports.UsersStorage
	verifyLinkStorage ports.EmailVerifyStorage
	publisher         publisher.Publisher
	query             queries.Users
}

func NewCreateUserUsecase(
	txm pgxtx.Manager,
	log *slog.Logger,
	cfg *EmailConfig,
	usersStorage ports.UsersStorage,
	verifyLinkStorage ports.EmailVerifyStorage,
	publisher publisher.Publisher,
	query queries.Users,
) *CreateUserUsecase {
	return &CreateUserUsecase{
		txm:               txm,
		log:               log,
		cfg:               cfg,
		usersStorage:      usersStorage,
		verifyLinkStorage: verifyLinkStorage,
		publisher:         publisher,
		query:             query,
	}
}

func (u *CreateUserUsecase) Execute(ctx context.Context, in *inputs.CreateUser, locale string) (string, error) {

	const op = "usecases.create_user"

	var createUserResult *domain.UserPrimaryFields

	err := u.txm.Wrap(ctx, func(txCtx context.Context) error {

		var err error

		_, err = u.query.FindByEmail(txCtx, in.GetEmail())

		if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
			return fmt.Errorf("%s: %w", op, err)
		}

		if err == nil {
			return fmt.Errorf("%s: %w", op, errs.ErrUserAlreadyExists)
		}

		hashedPassword, err := utils.HashPassword(in.GetPassword())

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		in.SetPassword(hashedPassword)

		createUserResult, err = u.usersStorage.Create(txCtx, in)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		expiresAt := time.Now().Add(u.cfg.EmailVerifyLinkExpiresDuration).UTC()

		link, err := u.verifyLinkStorage.Create(txCtx, createUserResult.ID, expiresAt)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.publisher.UserCreated(txCtx, &publisher.UserCreatedParams{
			ID:     createUserResult.ID,
			Email:  createUserResult.Email,
			Locale: locale,
		})

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.publisher.UserEmalVerifyLinkAdded(txCtx, &publisher.UserEmalVerifyLinkAddedParams{
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
