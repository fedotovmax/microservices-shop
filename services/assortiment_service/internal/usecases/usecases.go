package usecases

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/inputs"
	"github.com/fedotovmax/pgxtx"
)

type Config struct{}

type BrandsStorage interface {
	Create(ctx context.Context, in *inputs.CreateBrand, slug string) error
	FindBy(ctx context.Context, column db.BrandEntityFields, searchValue string) (*domain.Brand, error)
	GetAll(ctx context.Context) ([]domain.Brand, error)
	Update(ctx context.Context, params *db.UpdateBrandParams) error
	Delete(ctx context.Context, column db.BrandEntityFields, searchTerm string) error
}

type usecases struct {
	txm           pgxtx.Manager
	log           *slog.Logger
	cfg           *Config
	brandsStorage BrandsStorage
}

func New(
	txm pgxtx.Manager,
	log *slog.Logger,
	cfg *Config,
	brandsStorage BrandsStorage,
) *usecases {
	return &usecases{
		txm:           txm,
		log:           log,
		cfg:           cfg,
		brandsStorage: brandsStorage,
	}
}
