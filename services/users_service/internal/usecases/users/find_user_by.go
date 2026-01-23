package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter/db"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
)

func (u *usecases) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {

	const op = "usecases.FindUserByEmail"

	user, err := u.storage.FindUserBy(ctx, db.UserFieldEmail, email)

	if err != nil {
		if errors.Is(err, adapter.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w: %v", op, errs.ErrUserNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *usecases) FindUserByID(ctx context.Context, id string) (*domain.User, error) {
	const op = "usecases.users.FindUserByID"

	user, err := u.storage.FindUserBy(ctx, db.UserFieldID, id)

	if err != nil {
		if errors.Is(err, adapter.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w: %v", op, errs.ErrUserNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
