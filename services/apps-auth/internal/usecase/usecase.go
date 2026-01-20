package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop/apps-auth/internal/domain"
)

type Storage interface {
	SaveApp(ctx context.Context, secretHash string, app *domain.App) error
	FindApp(ctx context.Context, secretHash string) (*domain.App, error)
}

type Config struct {
	TokenSecret      string
	Issuer           string
	TokenExpDuration time.Duration
}

type usecases struct {
	log     *slog.Logger
	storage Storage
	cfg     *Config
}

func NewUsecases(log *slog.Logger, storage Storage, cfg *Config) *usecases {
	return &usecases{
		log:     log,
		storage: storage,
		cfg:     cfg,
	}
}
