package security

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
)

func (u *usecases) checkActiveBypass(ctx context.Context, user *domain.SessionsUser, code string) error {

	const op = "usecases.security.checkActiveBypass"

	if user.Bypass.IsCodeExpired() {
		err := u.AddLoginIPBypass(ctx, user)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, errs.ErrBypassCodeExpired)
	}

	if !user.Bypass.ComapreCodes(code) {
		return fmt.Errorf("%s: %w", op, errs.ErrBadBypassCode)
	}

	err := u.storage.RemoveSecurityBlock(ctx, db.SecurityTableBypass, user.Info.UID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
