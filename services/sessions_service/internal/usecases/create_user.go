package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/ports"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/queries"
)

type CreateUserUsecase struct {
	log          *slog.Logger
	usersStorage ports.UsersStorage
	query        queries.User
}

func NewCreateUserUsecase(
	log *slog.Logger,
	usersStorage ports.UsersStorage,
	query queries.User,
) *CreateUserUsecase {
	return &CreateUserUsecase{
		log:          log,
		usersStorage: usersStorage,
		query:        query,
	}
}

func (u *CreateUserUsecase) Execute(ctx context.Context, uid string, email string) error {
	const op = "usecases.create_user"

	_, err := u.query.FindByID(ctx, uid)

	if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err == nil {
		return fmt.Errorf("%s: %w", op, errs.ErrUserAlreadyExists)
	}

	err = u.usersStorage.Create(ctx, uid, email)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, errs.ErrInternalCreateUser, err)
	}

	return nil

}
