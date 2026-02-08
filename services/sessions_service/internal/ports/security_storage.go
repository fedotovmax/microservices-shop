package ports

import (
	"context"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
)

type SecurityStorage interface {
	RevokeTrustTokens(ctx context.Context, hashes []string) error
	UpdateTrustToken(ctx context.Context, in *inputs.CreateTrustToken) error
	FindUserTrustTokens(ctx context.Context, uid string) ([]*domain.DeviceTrustToken, error)
	FindTrustToken(ctx context.Context, uid, tokenHash string) (*domain.DeviceTrustToken, error)
	CreateTrustToken(ctx context.Context, in *inputs.CreateTrustToken) error

	AddSecurityBlock(ctx context.Context, operation db.Operation, table db.SecurityTable, in *inputs.Security) error
	RemoveSecurityBlock(ctx context.Context, table db.SecurityTable, uid string) error
}
