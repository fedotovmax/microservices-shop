package brands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/ports"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/queries"
	"github.com/gosimple/slug"
)

type CreateBrandUsecase struct {
	log     *slog.Logger
	storage ports.BrandsStorage
	query   queries.Brand
}

func NewCreateBrandUsecase(
	log *slog.Logger,
	storage ports.BrandsStorage,
	query queries.Brand,
) *CreateBrandUsecase {
	return &CreateBrandUsecase{
		log:     log,
		storage: storage,
		query:   query,
	}
}

func (u *CreateBrandUsecase) Execute(
	ctx context.Context,
	in *inputs.CreateBrand,
) error {

	const op = "usecases.brands.create"

	brandSlug := slug.Make(in.Title)

	_, err := u.query.FindBySlug(ctx, brandSlug)

	if err != nil && !errors.Is(err, errs.ErrBrandNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err == nil {
		return fmt.Errorf("%s: %w", op, errs.ErrBrandAlreadyExists)
	}

	err = u.storage.Create(ctx, in, brandSlug)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
