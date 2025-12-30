package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/pkg/logger"
)

func (u *usecases) addLoginIPBypassFn(ctx context.Context, l *slog.Logger, user *domain.SessionsUser) (*domain.SessionsUser, error) {

	const op = "usecases.addToBlacklistFn"

	var err error

	code, err := u.generateSecurityCode(u.cfg.LoginBypassCodeLength)

	if err != nil {
		l.Error("error when generate code for bypass", slog.String("uid", user.Info.UID), logger.Err(err))
		return nil, err
	}

	codeExpiresAt := time.Now().Add(u.cfg.LoginBypassExpDuration)

	bypassInput := &inputs.SecurityInput{
		UID:           user.Info.UID,
		Code:          code,
		CodeExpiresAt: codeExpiresAt,
	}

	if user.HasBypass() {
		err = u.storage.UpdateIPBypass(ctx, bypassInput)
	} else {
		err = u.storage.AddIPBypass(ctx, bypassInput)
	}

	if err != nil {
		l.Error("error when add/update bypass", slog.String("uid", user.Info.UID), logger.Err(err))
		return nil, err
	}

	updatedUser := user.Clone()

	updatedUser.Bypass = &domain.Bypass{
		Code:            code,
		BypassExpiresAt: codeExpiresAt,
	}

	//TODO: send event to kafka!

	return &updatedUser, nil
}

func (u *usecases) AddLoginIPBypass(ctx context.Context, inTx bool, user *domain.SessionsUser) (*domain.SessionsUser, error) {

	const op = "usecases.LoginBypass"

	l := u.log.With(slog.String("op", op))

	var updatedUser *domain.SessionsUser

	var err error

	if inTx {
		updatedUser, err = u.addLoginIPBypassFn(ctx, l, user)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return updatedUser, nil
	}

	err = u.txm.Wrap(ctx, func(txCtx context.Context) error {
		updatedUser, err = u.addLoginIPBypassFn(txCtx, l, user)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return updatedUser, nil

}
