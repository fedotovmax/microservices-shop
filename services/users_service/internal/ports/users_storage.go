package ports

import (
	"context"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
)

type UsersStorage interface {
	Create(ctx context.Context, d *inputs.CreateUser) (*domain.UserPrimaryFields, error)
	UpdateProfile(ctx context.Context, id string, in *inputs.UpdateUser) error
	FindBy(ctx context.Context, column db.UserEntityFields, value string) (*domain.User, error)
	SetIsEmailVerified(ctx context.Context, uid string, flag bool) error
}
