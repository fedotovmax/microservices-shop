package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter/db"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/pgxtx"
)

type UsersStorage interface {
	CreateUser(ctx context.Context, d *inputs.CreateUserInput) (*domain.UserPrimaryFields, error)
	UpdateUserProfile(ctx context.Context, id string, in *inputs.UpdateUserInput) error
	FindUserBy(ctx context.Context, column db.UserEntityFields, value string) (*domain.User, error)

	CreateEmailVerifyLink(ctx context.Context, userID string, expiresAt time.Time) (*domain.EmailVerifyLink, error)
	FindEmailVerifyLink(ctx context.Context, link string) (*domain.EmailVerifyLink, error)
	UpdateEmailVerifyLinkByUserID(ctx context.Context, userID string) (*domain.EmailVerifyLink, error)
}

type EventsStorage interface {
	SetEventStatusDone(ctx context.Context, id string) error
	SetEventsReservedToByIDs(ctx context.Context, ids []string, dur time.Duration) error
	RemoveEventReserve(ctx context.Context, id string) error
	CreateEvent(ctx context.Context, d *inputs.CreateEvent) (string, error)
	FindNewAndNotReservedEvents(ctx context.Context, limit int) ([]*domain.Event, error)
}

type Storage struct {
	users  UsersStorage
	events EventsStorage
}

type usecases struct {
	s   *Storage
	txm pgxtx.Manager
	log *slog.Logger
}

func CreateStorage(events EventsStorage, users UsersStorage) *Storage {
	return &Storage{
		events: events,
		users:  users,
	}
}

func NewUsecases(s *Storage, txm pgxtx.Manager, log *slog.Logger) *usecases {
	return &usecases{
		s:   s,
		txm: txm,
		log: log,
	}
}
