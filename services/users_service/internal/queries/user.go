package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/ports"
)

type Users interface {
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type users struct {
	usersStorage ports.UsersStorage
}

func NewUsers(usersStorage ports.UsersStorage) Users {
	return &users{
		usersStorage: usersStorage,
	}
}

func (q *users) FindByEmail(ctx context.Context, email string) (*domain.User, error) {

	user, err := q.usersStorage.FindBy(ctx, db.UserFieldEmail, email)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
		}
		return nil, err
	}

	return user, nil
}

func (q *users) FindByID(ctx context.Context, id string) (*domain.User, error) {

	user, err := q.usersStorage.FindBy(ctx, db.UserFieldID, id)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
		}
		return nil, err
	}

	return user, nil
}
