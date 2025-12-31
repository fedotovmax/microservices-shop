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

func (u *usecases) AddToBlacklist(ctx context.Context, user *domain.SessionsUser) (*domain.SessionsUser, error) {

	const op = "usecases.AddToBlackList"

	l := u.log.With(slog.String("op", op))

	var err error

	code, err := u.generateSecurityCode(u.cfg.BlacklistCodeLength)

	if err != nil {
		l.Error("error when generate code for blacklist", slog.String("uid", user.Info.UID), logger.Err(err))
		return nil, err
	}

	codeExpiresAt := time.Now().Add(u.cfg.BlacklistCodeExpDuration)

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
		return nil, err
	}

	updatedUser := user.Clone()

	updatedUser.BlackList = &domain.BlackList{
		Code:          code,
		CodeExpiresAt: codeExpiresAt,
	}

	eventPayload := events.SessionBlacklistAddedEventPayload{
		UID:           updatedUser.Info.UID,
		Email:         updatedUser.Info.Email,
		Code:          updatedUser.BlackList.Code,
		CodeExpiresAt: updatedUser.BlackList.CodeExpiresAt,
	}

	eventPayloadBytes, err := json.Marshal(eventPayload)

	if err != nil {
		return nil, err
	}

	eventInput := inputs.NewCreateEventInput()
	eventInput.SetAggregateID(updatedUser.Info.UID)
	eventInput.SetTopic(events.SESSION_EVENTS)
	eventInput.SetType(events.SESSION_BLACKLIST_ADDED)
	eventInput.SetPayload(eventPayloadBytes)

	_, err = u.storage.CreateEvent(ctx, eventInput)

	if err != nil {
		return nil, err
	}

	return &updatedUser, nil

}
