package ports

import (
	"context"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain"
)

type LanguagesStorage interface {
	Update(ctx context.Context, code string, isDefault, isActive bool) error
	GetDefault(ctx context.Context) (*domain.Language, error)
	FindAll(ctx context.Context) ([]domain.Language, error)
	Delete(ctx context.Context, code string) error
	Add(ctx context.Context, code string, isDefault, isActive bool) error
}
