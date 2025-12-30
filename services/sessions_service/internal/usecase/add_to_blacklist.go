package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/pkg/logger"
)

func (u *usecases) addToBlacklistFn(ctx context.Context, l *slog.Logger, session *domain.Session) (*domain.Session, error) {

	const op = "usecases.addToBlacklistFn"

	var err error

	code, err := u.generateSecurityCode(u.cfg.BlacklistCodeLength)

	if err != nil {
		l.Error("error when generate code for blacklist", slog.String("sid", session.ID), slog.String("uid", session.User.Info.UID), logger.Err(err))
		return nil, err
	}

	codeExpiresAt := time.Now().Add(u.cfg.BlacklistCodeExpDuration)

	blacklistInput := &inputs.SecurityInput{
		UID:           session.User.Info.UID,
		Code:          code,
		CodeExpiresAt: codeExpiresAt,
	}

	if session.User.IsInBlackList() {
		err = u.storage.UpdateBlacklistCode(ctx, blacklistInput)
	} else {
		err = u.storage.AddToBlackList(ctx, blacklistInput)
	}

	if err != nil {
		l.Error("error when add/update blacklist", slog.String("sid", session.ID), slog.String("uid", session.User.Info.UID), logger.Err(err))
		return nil, err
	}

	updatedSession := session.Clone()

	updatedSession.User.BlackList = &domain.BlackList{
		Code:          code,
		CodeExpiresAt: codeExpiresAt,
	}

	eventPayload := events.SessionBlacklistAddedEventPayload{
		UID:           updatedSession.User.Info.UID,
		Email:         updatedSession.User.Info.Email,
		Code:          updatedSession.User.BlackList.Code,
		CodeExpiresAt: updatedSession.User.BlackList.CodeExpiresAt,
	}

	eventPayloadBytes, err := json.Marshal(eventPayload)

	if err != nil {
		return nil, err
	}

	eventInput := inputs.NewCreateEventInput()
	eventInput.SetAggregateID(session.User.Info.UID)
	eventInput.SetTopic(events.SESSION_EVENTS)
	eventInput.SetType(events.SESSION_BLACKLIST_ADDED)
	eventInput.SetPayload(eventPayloadBytes)

	_, err = u.storage.CreateEvent(ctx, eventInput)

	if err != nil {
		return nil, err
	}

	return &updatedSession, nil
}

func (u *usecases) AddToBlackList(ctx context.Context, inTx bool, session *domain.Session) (*domain.Session, error) {

	const op = "usecases.AddToBlackList"

	l := u.log.With(slog.String("op", op))

	var updatedSession *domain.Session
	var err error

	if inTx {
		updatedSession, err = u.addToBlacklistFn(ctx, l, session)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return updatedSession, nil
	}

	err = u.txm.Wrap(ctx, func(txCtx context.Context) error {
		updatedSession, err = u.addToBlacklistFn(txCtx, l, session)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return updatedSession, nil

}
