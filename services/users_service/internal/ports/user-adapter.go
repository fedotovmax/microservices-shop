package ports

import (
	"context"

	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
)

type UserAdapter interface {
	Create(ctx context.Context, d domain.CreateUser) (string, error)
}
