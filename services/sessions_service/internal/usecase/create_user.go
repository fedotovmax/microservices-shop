package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
)

func (u *usecases) CreateUser(ctx context.Context, uid string, email string) error {

	const op = "usecases.CreateUser"

	_, err := u.FindUserByID(ctx, uid)

	if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err == nil {
		return fmt.Errorf("%s: %w", op, errs.ErrUserAlreadyExists)
	}

	err = u.storage.sessions.CreateUser(ctx, uid, email)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, errs.ErrInternalCreateUser, err)
	}

	return nil
}
