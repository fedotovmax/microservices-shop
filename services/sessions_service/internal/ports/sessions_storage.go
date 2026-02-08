package ports

import (
	"context"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
)

type SessionsStorage interface {
	Create(ctx context.Context, in *inputs.CreateSession) (string, error)
	Revoke(ctx context.Context, sids []string) error
	FindBy(ctx context.Context, column db.SessionEntityFields, value string) (*domain.Session, error)
	Update(ctx context.Context, in *inputs.CreateSession) error
	FindAllByUserID(ctx context.Context, uid string) ([]*domain.Session, error)
}
