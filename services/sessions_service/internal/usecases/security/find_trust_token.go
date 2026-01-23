package security

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
)

func (u *usecases) findTrustToken(ctx context.Context, uid, hash string) (*domain.DeviceTrustToken, error) {

	const op = "usecases.security.findTrustToken"

	t, err := u.storage.FindTrustToken(ctx, uid, hash)

	if err != nil {
		if errors.Is(err, adapter.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w: %v", op, errs.ErrTrustTokenNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}
