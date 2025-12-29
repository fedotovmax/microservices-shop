package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
)

func (u *usecases) findSession(ctx context.Context, column db.SessionEntityFields, value string) (*domain.Session, error) {

	//TODO: when find session and if revoked will report this (leaked session)

	const op = "usecases.findSession"

	session, err := u.storage.FindSession(ctx, column, value)

	if err != nil {
		if errors.Is(err, adapter.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w: %v", op, errs.ErrSessionNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if session.User.IsInBlackList() {
		return nil, fmt.Errorf("%s: %w", op, errs.NewUserSessionsInBlacklistError(
			session.User.Info.Email,
			session.User.Info.UID,
		))
	}

	if session.IsRevoked() {
		return nil, fmt.Errorf("%s: %w", op, errs.NewUserSessionRevokedError(
			session.User.Info.Email,
			session.User.Info.UID,
			session.ID,
		))
	}

	if session.IsExpired() {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrSessionExpired)
	}

	return session, nil
}

func (u *usecases) FindSessionByHash(ctx context.Context, hash string) (*domain.Session, error) {

	const op = "usecases.FindSessionByHash"

	session, err := u.findSession(ctx, db.SessionFieldRefreshHash, hash)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil

}

func (u *usecases) FindSessionByID(ctx context.Context, sid string) (*domain.Session, error) {

	const op = "usecases.FindSessionByID"

	session, err := u.findSession(ctx, db.SessionFieldID, sid)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}
