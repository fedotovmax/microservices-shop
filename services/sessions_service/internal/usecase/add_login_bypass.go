package usecase

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/pkg/logger"
)

func (u *usecases) AddLoginIPBypass(ctx context.Context, user *domain.SessionsUser) (*domain.SessionsUser, error) {

	const op = "usecases.AddLoginIPBypass"

	l := u.log.With(slog.String("op", op))

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
		err = u.storage.AddSecurityBlock(ctx, db.OperationUpdate, db.SecurityTableBypass, bypassInput)
	} else {
		err = u.storage.AddSecurityBlock(ctx, db.OperationInsert, db.SecurityTableBypass, bypassInput)
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

	eventPayload := events.SessionBypassAddedEventPayload{
		UID:             updatedUser.Info.UID,
		Email:           updatedUser.Info.Email,
		Code:            updatedUser.Bypass.Code,
		BypassExpiresAt: updatedUser.Bypass.BypassExpiresAt,
	}

	eventPayloadBytes, err := json.Marshal(eventPayload)

	if err != nil {
		return nil, err
	}

	eventInput := inputs.NewCreateEventInput()
	eventInput.SetAggregateID(updatedUser.Info.UID)
	eventInput.SetTopic(events.SESSION_EVENTS)
	eventInput.SetType(events.SESSION_BYPASS_ADDED)
	eventInput.SetPayload(eventPayloadBytes)

	_, err = u.storage.CreateEvent(ctx, eventInput)

	if err != nil {
		return nil, err
	}

	return &updatedUser, nil

}
