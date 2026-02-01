package usecases

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
)

func (u *usecases) handleSessionRevoked(ctx context.Context, session *domain.Session) error {

	const op = "usecases.security.handleSessionRevoked"

	if session.IsRevoked() && !session.User.IsInBlackList() {

		err := u.AddToBlacklist(ctx, session.User)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return fmt.Errorf("%s: %w", op, errs.ErrUserSessionsInBlackList)
	}
	return nil
}
