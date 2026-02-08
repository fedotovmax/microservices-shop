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
	eventspublisher "github.com/fedotovmax/microservices-shop/sessions_service/internal/events_publisher"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/ports"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/utils"
	"github.com/fedotovmax/microservices-shop/sessions_service/pkg/logger"
)

type AddToBlacklistUsecase struct {
	log             *slog.Logger
	cfg             *SecurityConfig
	securityStorage ports.SecurityStorage
	publisher       eventspublisher.Publisher
}

func NewAddToBlacklistUsecase(
	log *slog.Logger,
	cfg *SecurityConfig,
	securityStorage ports.SecurityStorage,
	publisher eventspublisher.Publisher,
) *AddToBlacklistUsecase {
	return &AddToBlacklistUsecase{
		cfg:             cfg,
		log:             log,
		securityStorage: securityStorage,
		publisher:       publisher,
	}
}

func (u *AddToBlacklistUsecase) Execute(ctx context.Context, user *domain.SessionsUser) error {
	const op = "usecases.add_to_blacklist"

	l := u.log.With(slog.String("op", op))

	var err error

	code, err := utils.GenerateSecurityCode(u.cfg.BlacklistCodeLength)

	if err != nil {
		l.Error("error when generate code for blacklist", slog.String("uid", user.Info.UID), logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	codeExpiresAt := time.Now().Add(u.cfg.BlacklistCodeExpDuration).UTC()

	blacklistInput := &inputs.Security{
		UID:           user.Info.UID,
		Code:          code,
		CodeExpiresAt: codeExpiresAt,
	}

	if user.IsInBlackList() {
		err = u.securityStorage.AddSecurityBlock(ctx, db.OperationUpdate, db.SecurityTableBlacklist, blacklistInput)
	} else {
		err = u.securityStorage.AddSecurityBlock(ctx, db.OperationInsert, db.SecurityTableBlacklist, blacklistInput)
	}

	if err != nil {
		l.Error("error when add/update blacklist", slog.String("uid", user.Info.UID), logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.publisher.SessionBlacklistAdded(ctx, events.SessionBlacklistAddedEventPayload{
		UID:           user.Info.UID,
		Email:         user.Info.Email,
		Code:          blacklistInput.Code,
		CodeExpiresAt: blacklistInput.CodeExpiresAt,
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
