package security

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/utils"
	"github.com/fedotovmax/microservices-shop/sessions_service/pkg/logger"
)

func (u *usecases) AddLoginIPBypass(ctx context.Context, user *domain.SessionsUser) (*time.Time, error) {

	const op = "usecases.security.AddLoginIPBypass"

	l := u.log.With(slog.String("op", op))

	var err error

	code, err := utils.GenerateSecurityCode(u.cfg.LoginBypassCodeLength)

	if err != nil {
		l.Error("error when generate code for bypass", slog.String("uid", user.Info.UID), logger.Err(err))
		return nil, err
	}

	codeExpiresAt := time.Now().Add(u.cfg.LoginBypassExpDuration).UTC()

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

	eventPayload := events.SessionBypassAddedEventPayload{
		UID:             user.Info.UID,
		Email:           user.Info.Email,
		Code:            bypassInput.Code,
		BypassExpiresAt: bypassInput.CodeExpiresAt,
	}

	eventPayloadBytes, err := json.Marshal(eventPayload)

	if err != nil {
		return nil, err
	}

	eventInput := inputs.NewCreateEventInput()
	eventInput.SetAggregateID(user.Info.UID)
	eventInput.SetTopic(events.SESSION_EVENTS)
	eventInput.SetType(events.SESSION_BYPASS_ADDED)
	eventInput.SetPayload(eventPayloadBytes)

	_, err = u.eventSender.CreateEvent(ctx, eventInput)

	if err != nil {
		return nil, err
	}

	return &codeExpiresAt, nil

}
