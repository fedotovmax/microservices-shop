package usecases

import (
	"context"
	"log/slog"
	"time"

	"github.com/fedotovmax/kafka-lib/outbox"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapters/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/pgxtx"
	"github.com/medama-io/go-useragent"
)

type SessionsStorage interface {
	CreateSession(ctx context.Context, in *inputs.CreateSessionInput) (string, error)
	RevokeSessions(ctx context.Context, sids []string) error
	FindSession(ctx context.Context, column db.SessionEntityFields, value string) (*domain.Session, error)
	UpdateSession(ctx context.Context, in *inputs.CreateSessionInput) error

	FindUser(ctx context.Context, uid string) (*domain.SessionsUser, error)
	FindUserSessions(ctx context.Context, uid string) ([]*domain.Session, error)
	CreateUser(ctx context.Context, uid string, email string) error
}

type SecurityStorage interface {
	RevokeTrustTokens(ctx context.Context, hashes []string) error
	UpdateTrustToken(ctx context.Context, in *inputs.CreateTrustTokenInput) error
	FindUserTrustTokens(ctx context.Context, uid string) ([]*domain.DeviceTrustToken, error)
	FindTrustToken(ctx context.Context, uid, tokenHash string) (*domain.DeviceTrustToken, error)
	CreateTrustToken(ctx context.Context, in *inputs.CreateTrustTokenInput) error

	AddSecurityBlock(ctx context.Context, operation db.Operation, table db.SecurityTable, in *inputs.SecurityInput) error
	RemoveSecurityBlock(ctx context.Context, table db.SecurityTable, uid string) error
}

type EventSender interface {
	CreateEvent(ctx context.Context, d *outbox.CreateEvent) (string, error)
}

type Config struct {
	TokenIssuer string
	TokenSecret string

	RefreshExpiresDuration time.Duration

	AccessExpiresDuration time.Duration

	BlacklistCodeExpDuration time.Duration

	LoginBypassExpDuration time.Duration

	DeviceTrustTokenExpDuration time.Duration
	DeviceTrustTokenThreshold   time.Duration

	BlacklistCodeLength uint8

	LoginBypassCodeLength uint8
}

type usecases struct {
	sessionsStorage SessionsStorage
	securityStorage SecurityStorage
	eventSender     EventSender
	uaparser        *useragent.Parser
	txm             pgxtx.Manager
	cfg             *Config
	log             *slog.Logger
}

func New(
	sessionsStorage SessionsStorage,
	securityStorage SecurityStorage,
	eventSender EventSender,
	txm pgxtx.Manager,
	log *slog.Logger,
	cfg *Config,
) *usecases {
	return &usecases{
		sessionsStorage: sessionsStorage,
		securityStorage: securityStorage,
		eventSender:     eventSender,
		uaparser:        useragent.NewParser(),
		txm:             txm,
		cfg:             cfg,
		log:             log,
	}
}
