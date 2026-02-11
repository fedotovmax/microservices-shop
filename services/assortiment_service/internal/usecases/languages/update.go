package languages

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/ports"
)

type UpdateLanguageUsecase struct {
	log     *slog.Logger
	storage ports.LanguagesStorage
}

func NewUpdateLanguageUsecase(
	log *slog.Logger,
	storage ports.LanguagesStorage,
) *UpdateLanguageUsecase {
	return &UpdateLanguageUsecase{
		log:     log,
		storage: storage,
	}
}

func (u *UpdateLanguageUsecase) Execute(
	ctx context.Context,
	code string,
	isDefault bool,
	isActive bool,
) error {

	const op = "usecases.languages.update"

	err := u.storage.Update(ctx, code, isDefault, isActive)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
