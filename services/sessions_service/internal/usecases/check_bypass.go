package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/ports"
)

type CheckBypassUsecase struct {
	log             *slog.Logger
	securityStorage ports.SecurityStorage
	addLoginBypass  *AddLoginBypassUsecase
}

func NewCheckBypassUsecase(
	log *slog.Logger,
	securityStorage ports.SecurityStorage,
	addLoginBypass *AddLoginBypassUsecase,
) *CheckBypassUsecase {
	return &CheckBypassUsecase{
		log:             log,
		securityStorage: securityStorage,
		addLoginBypass:  addLoginBypass,
	}
}

func (u *CheckBypassUsecase) Execute(ctx context.Context, user *domain.SessionsUser, code string) error {

	const op = "usecases.check_bypass"

	if user.Bypass.IsCodeExpired() {
		//TODO:change this maybe?? and send new code on demand
		_, err := u.addLoginBypass.Execute(ctx, user)
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
