package brands

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/inputs"
	"github.com/fedotovmax/pgxtx"
	"github.com/gosimple/slug"
)

type CreateBrandUsecase struct {
	txm     pgxtx.Manager
	log     *slog.Logger
	storage BrandsStorage
}

func NewCreateBrandUsecase(
	txm pgxtx.Manager,
	log *slog.Logger,
	storage BrandsStorage,
) *CreateBrandUsecase {
	return &CreateBrandUsecase{
		txm:     txm,
		log:     log,
		storage: storage,
	}
}

func (uc *CreateBrandUsecase) Execute(
	ctx context.Context,
	in *inputs.CreateBrand,
) error {

	const op = "usecases.brands.create.Execute"

	//validatation

	brandslug := slug.Make(in.Title)

	err := uc.storage.Create(ctx, in, brandslug)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
