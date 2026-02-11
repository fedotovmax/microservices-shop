package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/assortiment_service/internal/ports"
)

type Brand interface {
	FindByID(ctx context.Context, id string) (*domain.Brand, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Brand, error)
}

type brand struct {
	storage ports.BrandsStorage
}

func NewBrand(s ports.BrandsStorage) Brand {
	return &brand{storage: s}
}

func (q *brand) FindByID(ctx context.Context, id string) (*domain.Brand, error) {
	b, err := q.storage.FindBy(ctx, db.BrandFieldID, id)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrBrandNotFound, err)
		}
		return nil, err
	}

	return b, nil
}

func (q *brand) FindBySlug(ctx context.Context, slug string) (*domain.Brand, error) {

	b, err := q.storage.FindBy(ctx, db.BrandFieldSlug, slug)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrBrandNotFound, err)
		}
		return nil, err
	}

	return b, nil
}

func (q *brand) FindAll(ctx context.Context, onlyActive bool) ([]domain.Brand, error) {

	b, err := q.storage.FindAll(ctx, onlyActive)

	if err != nil {
		return nil, err
	}

	return b, nil
}
