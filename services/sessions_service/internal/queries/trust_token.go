package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/ports"
)

type TrustToken interface {
	Find(ctx context.Context, uid, hash string) (*domain.DeviceTrustToken, error)
}

type trustToken struct {
	securityStorage ports.SecurityStorage
}

func NewTrustToken(securityStorage ports.SecurityStorage) TrustToken {
	return &trustToken{
		securityStorage: securityStorage,
	}
}

func (q *trustToken) Find(ctx context.Context, uid, hash string) (*domain.DeviceTrustToken, error) {

	t, err := q.securityStorage.FindTrustToken(ctx, uid, hash)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrTrustTokenNotFound, err)
		}
		return nil, fmt.Errorf("%w", err)
	}

	return t, nil
}
