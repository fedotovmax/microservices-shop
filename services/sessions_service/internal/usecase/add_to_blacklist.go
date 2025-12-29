package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/pkg/utils"
)

func (u *usecases) AddToBlackList(ctx context.Context, session *domain.Session) (*domain.Session, error) {

	const op = "usecases.AddToBlackList"

	l := u.log.With(slog.String("op", op))

	code, err := utils.GenerateCode(6)

	if err != nil {
		l.Error("error when generate code for blacklist", slog.String("sid", session.ID), slog.String("uid", session.User.Info.UID))
		return nil, fmt.Errorf("%s: %w: %v", op, errs.ErrSessionRevoked, err)
	}

	codeExpiresAt := time.Now().Add(6 * time.Hour)

	err = u.storage.AddToBlackList(ctx, &inputs.AddToBlackListInput{
		UID:           session.User.Info.UID,
		Code:          code,
		CodeExpiresAt: codeExpiresAt,
	})

	if err != nil {
		l.Error("error when add session to blacklist", slog.String("sid", session.ID), slog.String("uid", session.User.Info.UID))
		return nil, fmt.Errorf("%s: %w: %v", op, errs.ErrSessionRevoked, err)
	}

	session.User.BlackList = &domain.BlackList{
		Code:          code,
		CodeExpiresAt: codeExpiresAt,
	}

	return session, nil

}
