package languages

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/ports"
)

type AddLanguageUsecase struct {
	log     *slog.Logger
	storage ports.LanguagesStorage
}

func NewAddLanguageUsecase(
	log *slog.Logger,
	storage ports.LanguagesStorage,
) *AddLanguageUsecase {
	return &AddLanguageUsecase{
		log:     log,
		storage: storage,
	}
}

func (u *AddLanguageUsecase) Execute(
	ctx context.Context,
	code string,
	isDefault bool,
	isActive bool,
) error {

	const op = "usecases.languages.add"

	err := u.storage.Add(ctx, code, isDefault, isActive)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
