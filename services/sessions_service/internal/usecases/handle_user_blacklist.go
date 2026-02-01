package usecases

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
)

func (u *usecases) handleUserBlacklist(ctx context.Context, user *domain.SessionsUser) error {

	const op = "usecases.security.handleUserBlacklist"

	if user.IsInBlackList() {

		if user.BlackList.IsCodeExpired() {

			err := u.AddToBlacklist(ctx, user)

			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}

			return fmt.Errorf("%s: %w", op, errs.ErrBlacklistCodeExpired)
		}

		return fmt.Errorf("%s: %w", op, errs.ErrUserSessionsInBlackList)
	}

	return nil
}
