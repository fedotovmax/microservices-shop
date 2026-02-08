package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/ports"
)

type User interface {
	FindByID(ctx context.Context, uid string) (*domain.SessionsUser, error)
}

type user struct {
	usersStorage ports.UsersStorage
}

func NewUser(usersStorage ports.UsersStorage) User {
	return &user{
		usersStorage: usersStorage,
	}
}

func (q *user) FindByID(ctx context.Context, uid string) (*domain.SessionsUser, error) {

	user, err := q.usersStorage.Find(ctx, uid)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
		}
		return nil, fmt.Errorf("%w", err)
	}

	return user, nil
}
