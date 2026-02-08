package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
)

type IsUserInBlacklistUsecase struct {
	log            *slog.Logger
	addToBlacklist *AddToBlacklistUsecase
}

func NewIsUserInBlacklistUsecase(
	log *slog.Logger,
	addToBlacklist *AddToBlacklistUsecase,
) *IsUserInBlacklistUsecase {
	return &IsUserInBlacklistUsecase{
		log:            log,
		addToBlacklist: addToBlacklist,
	}
}

func (u *IsUserInBlacklistUsecase) Execute(ctx context.Context, user *domain.SessionsUser) error {

	const op = "usecases.is_user_in_blacklist"

	if user.IsInBlackList() {

		if user.BlackList.IsCodeExpired() {

			err := u.addToBlacklist.Execute(ctx, user)

			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}

			return fmt.Errorf("%s: %w", op, errs.ErrBlacklistCodeExpired)
		}

		return fmt.Errorf("%s: %w", op, errs.ErrUserSessionsInBlackList)
	}

	return nil
}
