package usecases

import (
	"context"
	"log/slog"
	"time"

	"github.com/fedotovmax/kafka-lib/outbox"
	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/pgxtx"
)

type UsersStorage interface {
	CreateUser(ctx context.Context, d *inputs.CreateUserInput) (*domain.UserPrimaryFields, error)
	UpdateUserProfile(ctx context.Context, id string, in *inputs.UpdateUserInput) error
	FindUserBy(ctx context.Context, column db.UserEntityFields, value string) (*domain.User, error)
	SetIsEmailVerified(ctx context.Context, uid string, flag bool) error
}

type EmailVerifyStorage interface {
	CreateEmailVerifyLink(ctx context.Context, userID string, expiresAt time.Time) (*domain.EmailVerifyLink, error)
	FindEmailVerifyLinkBy(ctx context.Context, column db.VerifyEmailLinkEntityFields, value string) (*domain.EmailVerifyLink, error)
	UpdateEmailVerifyLinkByUserID(ctx context.Context, userID string, expiresAt time.Time) (*domain.EmailVerifyLink, error)
	DeleteEmailVerifyLink(ctx context.Context, link string) error
}

type Config struct {
	EmailVerifyLinkExpiresDuration time.Duration
}

type EventSender interface {
	CreateEvent(ctx context.Context, d *outbox.CreateEvent) (string, error)
}

type usecases struct {
	usersStorage       UsersStorage
	emailVerifyStorage EmailVerifyStorage
	eventSender        EventSender
	txm                pgxtx.Manager
	log                *slog.Logger
	cfg                *Config
}

func New(
	usersStorage UsersStorage,
	emailVerifyStorage EmailVerifyStorage,
	txm pgxtx.Manager,
	es EventSender,
	log *slog.Logger,
	cfg *Config,
) *usecases {
	return &usecases{
		usersStorage:       usersStorage,
		emailVerifyStorage: emailVerifyStorage,
		eventSender:        es,
		txm:                txm,
		log:                log,
		cfg:                cfg,
	}
}
