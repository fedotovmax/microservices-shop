package brands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/ports"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/queries"
	"github.com/gosimple/slug"
)

type UpdateBrandUsecase struct {
	log     *slog.Logger
	storage ports.BrandsStorage
	query   queries.Brand
}

func NewUpdateBrandUsecase(
	log *slog.Logger,
	storage ports.BrandsStorage,
	query queries.Brand,
) *UpdateBrandUsecase {
	return &UpdateBrandUsecase{
		log:     log,
		storage: storage,
		query:   query,
	}
}

func (u *UpdateBrandUsecase) Execute(ctx context.Context, in *inputs.UpdateBrand) error {

	const op = "usecases.brands.update"

	var newSlug *string

	if in.Title != nil {

		brandSlug := slug.Make(*in.Title)

		_, err := u.query.FindBySlug(ctx, brandSlug)

		if err != nil && !errors.Is(err, errs.ErrBrandNotFound) {
			return fmt.Errorf("%s: %w", op, err)
		}

		if err == nil {
			return fmt.Errorf("%s: %w", op, errs.ErrBrandAlreadyExists)
		}

		newSlug = &brandSlug
	}

	err := u.storage.Update(ctx, &db.UpdateBrandParams{
		Input:        in,
		NewSlug:      newSlug,
		SearchColumn: db.BrandFieldID,
		SearchValue:  in.ID,
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
