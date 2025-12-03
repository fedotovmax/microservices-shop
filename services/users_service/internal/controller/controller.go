package controller

import (
	"context"

	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
)

type Usecases interface {
	CreateUser(ctx context.Context, d *domain.CreateUserInput) (string, error)
}
