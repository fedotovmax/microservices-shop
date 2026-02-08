package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/ports"
)

type RevokeSessionsUsecase struct {
	log             *slog.Logger
	sessionsStorage ports.SessionsStorage
}

func NewRevokeSessionsUsecase(
	log *slog.Logger,
	sessionsStorage ports.SessionsStorage,
) *RevokeSessionsUsecase {
	return &RevokeSessionsUsecase{
		log:             log,
		sessionsStorage: sessionsStorage,
	}
}

func (u *RevokeSessionsUsecase) Execute(ctx context.Context, sids []string) error {

	const op = "usecases.revoke_sessions"

	err := u.sessionsStorage.Revoke(ctx, sids)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
