package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/ports"
)

type Session interface {
	FindByHash(ctx context.Context, hash string) (*domain.Session, error)
	FindByID(ctx context.Context, id string) (*domain.Session, error)
}

type session struct {
	sessionsStorage ports.SessionsStorage
}

func NewSession(sessionsStorage ports.SessionsStorage) Session {
	return &session{
		sessionsStorage: sessionsStorage,
	}
}

func (q *session) FindByHash(ctx context.Context, hash string) (*domain.Session, error) {

	foundedSession, err := q.sessionsStorage.FindBy(ctx, db.SessionFieldRefreshHash, hash)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
		}
		return nil, err
	}

	return foundedSession, nil
}

func (q *session) FindByID(ctx context.Context, id string) (*domain.Session, error) {

	foundedSession, err := q.sessionsStorage.FindBy(ctx, db.SessionFieldID, id)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
		}
		return nil, err
	}

	return foundedSession, nil
}
