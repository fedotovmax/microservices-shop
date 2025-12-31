package usecase

import (
	"context"
	"fmt"

	"github.com/fedotovmax/goutils/sliceutils"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
)

func (u *usecases) handleUserBlacklist(ctx context.Context, user *domain.SessionsUser) error {

	const op = "usecases.handleUserBlacklist"

	if user.IsInBlackList() {

		if user.BlackList.IsCodeExpired() {

			_, err := u.AddToBlacklist(ctx, user)

			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}

			return fmt.Errorf("%s: %w", op, errs.ErrBlacklistCodeExpired)
		}

		return fmt.Errorf("%s: %w", op, errs.ErrUserSessionsInBlackList)
	}

	return nil
}

type bypassParams struct {
	IP         string
	Browser    string
	BypassCode string
}

func (u *usecases) handleUserBypass(ctx context.Context, user *domain.SessionsUser, params bypassParams) error {

	const op = "usecases.handleUserBypass"

	if !user.HasBypass() {
		userActiveSessions, err := u.storage.FindUserSessions(ctx, user.Info.UID)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		userActiveSessions = sliceutils.Filter(userActiveSessions, func(session *domain.Session) bool {
			return !session.IsRevoked() && !session.IsExpired()
		})

		_, trusted := sliceutils.Find(userActiveSessions, func(session *domain.Session) bool {
			return session.IP == params.IP && session.Browser == params.Browser
		})

		if !trusted {

			_, err := u.AddLoginIPBypass(ctx, user)

			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}

			return fmt.Errorf("%s: %w", op, errs.ErrLoginFromNewIPOrDevice)
		}

		return nil
	}

	if user.Bypass.IsCodeExpired() {

		_, err := u.AddLoginIPBypass(ctx, user)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return fmt.Errorf("%s: %w", op, errs.ErrBypassCodeExpired)
	}

	if !user.Bypass.ComapreCodes(params.BypassCode) {
		return fmt.Errorf("%s: %w", op, errs.ErrBadBypassCode)
	}

	err := u.storage.RemoveSecurityBlock(ctx, db.SecurityTableBypass, user.Info.UID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
