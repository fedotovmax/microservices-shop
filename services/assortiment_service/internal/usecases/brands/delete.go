package brands

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/ports"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/queries"
)

type DeleteBrandUsecase struct {
	log     *slog.Logger
	storage ports.BrandsStorage
	query   queries.Brand
}

func NewDeleteBrandUsecase(
	log *slog.Logger,
	storage ports.BrandsStorage,
	query queries.Brand,
) *DeleteBrandUsecase {
	return &DeleteBrandUsecase{
		log:     log,
		storage: storage,
		query:   query,
	}
}

func (u *DeleteBrandUsecase) Execute(ctx context.Context, id string) error {

	const op = "usecases.brands.delete"

	b, err := u.query.FindByID(ctx, id)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.storage.Delete(ctx, db.BrandFieldID, b.ID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
