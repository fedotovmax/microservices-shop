package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
)

type IsSessionRevokedUsecase struct {
	log            *slog.Logger
	addToBlacklist *AddToBlacklistUsecase
}

func NewIsSessionRevokedUsecase(
	log *slog.Logger,
	addToBlacklist *AddToBlacklistUsecase,
) *IsSessionRevokedUsecase {
	return &IsSessionRevokedUsecase{
		log:            log,
		addToBlacklist: addToBlacklist,
	}
}

func (u *IsSessionRevokedUsecase) Execute(ctx context.Context, session *domain.Session) error {
	const op = "usecases.is_session_revoked"

	if session.IsRevoked() && !session.User.IsInBlackList() {

		err := u.addToBlacklist.Execute(ctx, session.User)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return fmt.Errorf("%s: %w", op, errs.ErrUserSessionsInBlackList)
	}

	return nil
}
