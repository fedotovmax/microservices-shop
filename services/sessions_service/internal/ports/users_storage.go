package ports

import (
	"context"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
)

type UsersStorage interface {
	Find(ctx context.Context, uid string) (*domain.SessionsUser, error)
	Create(ctx context.Context, uid string, email string) error
}
