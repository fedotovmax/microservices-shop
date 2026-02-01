package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
)

func (u *usecases) FindUserByID(ctx context.Context, uid string) (*domain.SessionsUser, error) {

	const op = "usecases.security.FindUserByID"

	user, err := u.sessionsStorage.FindUser(ctx, uid)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w: %v", op, errs.ErrUserNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
