package brands

import (
	"context"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/inputs"
)

type BrandsStorage interface {
	Create(ctx context.Context, in *inputs.CreateBrand, slug string) error
	FindBy(ctx context.Context, column db.BrandEntityFields, searchValue string) (*domain.Brand, error)
	GetAll(ctx context.Context) ([]domain.Brand, error)
	Update(ctx context.Context, params *db.UpdateBrandParams) error
	Delete(ctx context.Context, column db.BrandEntityFields, searchTerm string) error
}
