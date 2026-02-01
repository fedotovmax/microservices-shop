package usecases

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
)

func (u *usecases) checkActiveBypass(ctx context.Context, user *domain.SessionsUser, code string) error {

	const op = "usecases.security.checkActiveBypass"

	if user.Bypass.IsCodeExpired() {
		//TODO:change this maybe?? and send new code on demand
		_, err := u.AddLoginIPBypass(ctx, user)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, errs.ErrBypassCodeExpired)
	}

	if !user.Bypass.ComapreCodes(code) {
		return fmt.Errorf("%s: %w", op, errs.ErrBadBypassCode)
	}

	err := u.securityStorage.RemoveSecurityBlock(ctx, db.SecurityTableBypass, user.Info.UID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
