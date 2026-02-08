package ports

import (
	"context"
	"time"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
)

type EmailVerifyStorage interface {
	Create(ctx context.Context, userID string, expiresAt time.Time) (*domain.EmailVerifyLink, error)
	FindBy(ctx context.Context, column db.VerifyEmailLinkEntityFields, value string) (*domain.EmailVerifyLink, error)
	UpdateByUserID(ctx context.Context, userID string, expiresAt time.Time) (*domain.EmailVerifyLink, error)
	Delete(ctx context.Context, link string) error
}
