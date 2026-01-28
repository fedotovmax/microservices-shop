package security

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/fedotovmax/kafka-lib/outbox"
	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/utils"
	"github.com/fedotovmax/microservices-shop/sessions_service/pkg/logger"
)

func (u *usecases) AddToBlacklist(ctx context.Context, user *domain.SessionsUser) error {

	const op = "usecases.security.AddToBlackList"

	l := u.log.With(slog.String("op", op))

	var err error

	code, err := utils.GenerateSecurityCode(u.cfg.BlacklistCodeLength)

	if err != nil {
		l.Error("error when generate code for blacklist", slog.String("uid", user.Info.UID), logger.Err(err))
		return err
	}

	codeExpiresAt := time.Now().Add(u.cfg.BlacklistCodeExpDuration).UTC()

	blacklistInput := &inputs.SecurityInput{
		UID:           user.Info.UID,
		Code:          code,
		CodeExpiresAt: codeExpiresAt,
	}

	if user.IsInBlackList() {
		err = u.storage.AddSecurityBlock(ctx, db.OperationUpdate, db.SecurityTableBlacklist, blacklistInput)
	} else {
		err = u.storage.AddSecurityBlock(ctx, db.OperationInsert, db.SecurityTableBlacklist, blacklistInput)
	}

	if err != nil {
		l.Error("error when add/update blacklist", slog.String("uid", user.Info.UID), logger.Err(err))
		return err
	}

	eventPayload := events.SessionBlacklistAddedEventPayload{
		UID:           user.Info.UID,
		Email:         user.Info.Email,
		Code:          blacklistInput.Code,
		CodeExpiresAt: blacklistInput.CodeExpiresAt,
	}

	eventPayloadBytes, err := json.Marshal(eventPayload)

	if err != nil {
		return err
	}

	eventInput := outbox.NewCreateEventInput()
	eventInput.SetAggregateID(user.Info.UID)
	eventInput.SetTopic(events.SESSION_EVENTS)
	eventInput.SetType(events.SESSION_BLACKLIST_ADDED)
	eventInput.SetPayload(eventPayloadBytes)

	_, err = u.eventSender.CreateEvent(ctx, eventInput)

	if err != nil {
		return err
	}

	return nil

}
