package languages

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/ports"
)

type DeleteLanguageUsecase struct {
	log     *slog.Logger
	storage ports.LanguagesStorage
}

func NewDeleteLanguageUsecase(
	log *slog.Logger,
	storage ports.LanguagesStorage,
) *DeleteLanguageUsecase {
	return &DeleteLanguageUsecase{
		log:     log,
		storage: storage,
	}
}

func (u *DeleteLanguageUsecase) Execute(
	ctx context.Context,
	code string,
) error {

	const op = "usecases.languages.delete"

	err := u.storage.Delete(ctx, code)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
