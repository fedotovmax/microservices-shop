package usecases

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/utils"
)

func (u *usecases) VerifyRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {

	const op = "usecases.security.VerifyRefreshToken"

	refreshTokenHash := utils.HashToken(refreshToken)

	session, err := u.FindSessionByHash(ctx, refreshTokenHash)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if session.User.IsDeleted() {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrUserDeleted)
	}

	if session.IsExpired() {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrSessionExpired)
	}

	err = u.handleUserBlacklist(ctx, session.User)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = u.handleSessionRevoked(ctx, session)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil

}
