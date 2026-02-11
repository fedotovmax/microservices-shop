package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/ports"
)

type Language interface {
	GetDefault(ctx context.Context) (*domain.Language, error)
}

type language struct {
	storage ports.LanguagesStorage
}

func NewLanguage(s ports.LanguagesStorage) Language {
	return &language{
		storage: s,
	}
}

func (q *language) GetDefault(ctx context.Context) (*domain.Language, error) {
	lang, err := q.storage.GetDefault(ctx)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrLanguageNotFound, err)
		}
		return nil, err
	}

	return lang, nil
}
