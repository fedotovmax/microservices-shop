package usecases

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/ports"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/publisher"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/utils"
	"github.com/fedotovmax/microservices-shop/sessions_service/pkg/logger"
)

type AddLoginBypassUsecase struct {
	log             *slog.Logger
	cfg             *SecurityConfig
	securityStorage ports.SecurityStorage
	publisher       publisher.Publisher
}

func NewAddLoginBypassUsecase(
	log *slog.Logger,
	cfg *SecurityConfig,
	securityStorage ports.SecurityStorage,
	publisher publisher.Publisher,
) *AddLoginBypassUsecase {
	return &AddLoginBypassUsecase{
		cfg:             cfg,
		log:             log,
		securityStorage: securityStorage,
		publisher:       publisher,
	}
}

func (u *AddLoginBypassUsecase) Execute(ctx context.Context, user *domain.SessionsUser) (*time.Time, error) {
	const op = "usecases.add_login_bypass"

	l := u.log.With(slog.String("op", op))

	var err error

	code, err := utils.GenerateSecurityCode(u.cfg.LoginBypassCodeLength)

	if err != nil {
		l.Error("error when generate code for bypass", slog.String("uid", user.Info.UID), logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	codeExpiresAt := time.Now().Add(u.cfg.LoginBypassExpDuration).UTC()

	bypassInput := &inputs.Security{
		UID:           user.Info.UID,
		Code:          code,
		CodeExpiresAt: codeExpiresAt,
	}

	if user.HasBypass() {
		err = u.securityStorage.AddSecurityBlock(ctx, db.OperationUpdate, db.SecurityTableBypass, bypassInput)
	} else {
		err = u.securityStorage.AddSecurityBlock(ctx, db.OperationInsert, db.SecurityTableBypass, bypassInput)
	}

	if err != nil {
		l.Error("error when add/update bypass", slog.String("uid", user.Info.UID), logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = u.publisher.SessionBypassAdded(ctx, events.SessionBypassAddedEventPayload{
		UID:             user.Info.UID,
		Email:           user.Info.Email,
		Code:            bypassInput.Code,
		BypassExpiresAt: bypassInput.CodeExpiresAt,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &codeExpiresAt, nil

}
