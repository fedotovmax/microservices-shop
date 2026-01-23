package users

import (
	"context"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter/db"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/pgxtx"
)

type Storage interface {
	CreateUser(ctx context.Context, d *inputs.CreateUserInput) (*domain.UserPrimaryFields, error)
	UpdateUserProfile(ctx context.Context, id string, in *inputs.UpdateUserInput) error
	FindUserBy(ctx context.Context, column db.UserEntityFields, value string) (*domain.User, error)

	CreateEmailVerifyLink(ctx context.Context, userID string, expiresAt time.Time) (*domain.EmailVerifyLink, error)
	FindEmailVerifyLink(ctx context.Context, link string) (*domain.EmailVerifyLink, error)
	UpdateEmailVerifyLinkByUserID(ctx context.Context, userID string) (*domain.EmailVerifyLink, error)
}

type Config struct {
	EmailVerifyLinkExpiresDuration time.Duration
}

type EventSender interface {
	CreateEvent(ctx context.Context, d *inputs.CreateEvent) (string, error)
}

type usecases struct {
	storage     Storage
	eventSender EventSender
	txm         pgxtx.Manager
	log         *slog.Logger
	cfg         *Config
}

func New(s Storage, txm pgxtx.Manager, es EventSender, log *slog.Logger, cfg *Config) *usecases {
	return &usecases{
		storage:     s,
		eventSender: es,
		txm:         txm,
		log:         log,
		cfg:         cfg,
	}
}
