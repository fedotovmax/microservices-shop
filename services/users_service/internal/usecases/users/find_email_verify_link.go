package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter/db"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
)

func (u *usecases) FindEmailVerifyLinkByUserID(ctx context.Context, uid string) (*domain.EmailVerifyLink, error) {

	const op = "FindEmailVerifyLinkByUserID"

	linkEntity, err := u.storage.FindEmailVerifyLinkBy(ctx, db.VerifyEmailLinkUserIDField, uid)

	if err != nil {
		if errors.Is(err, adapter.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w: %v", op, errs.ErrVerifyEmailLinkNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return linkEntity, nil
}

func (u *usecases) FindEmailVerifyLinkByPrimary(ctx context.Context, link string) (*domain.EmailVerifyLink, error) {

	const op = "FindEmailVerifyLinkByPrimary"

	linkEntity, err := u.storage.FindEmailVerifyLinkBy(ctx, db.VerifyEmailLinkPrimaryField, link)

	if err != nil {
		if errors.Is(err, adapter.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w: %v", op, errs.ErrVerifyEmailLinkNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return linkEntity, nil
}
