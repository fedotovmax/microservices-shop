package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter/db"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/pgxtx"
	"github.com/medama-io/go-useragent"
)

type JwtAdapter interface {
	CreateAccessToken(issuer, uid, sid string) (*domain.NewAccessToken, error)
	ParseAccessToken(token string, issuer string) (jti string, uid string, perr error)
}

type Storage interface {
	CreateSession(ctx context.Context, in *inputs.CreateSessionInput) (string, error)
	RevokeSessions(ctx context.Context, sids []string) error
	FindSession(ctx context.Context, column db.SessionEntityFields, value string) (*domain.Session, error)
	UpdateSession(ctx context.Context, in *inputs.CreateSessionInput) error
	AddToBlackList(ctx context.Context, in *inputs.AddToBlackListInput) error
	RemoveUserFromBlacklist(ctx context.Context, uid string) error
	FindUser(ctx context.Context, uid string) (*domain.SessionsUser, error)
	FindUserSessions(ctx context.Context, uid string) ([]*domain.Session, error)
}

type usecases struct {
	log                    *slog.Logger
	txm                    pgxtx.Manager
	jwt                    JwtAdapter
	storage                Storage
	uaparser               *useragent.Parser
	refreshExpiresDuration time.Duration
}

func New(log *slog.Logger, txm pgxtx.Manager, jwt JwtAdapter, storage Storage, refreshExpDuration time.Duration) *usecases {
	uaparser := useragent.NewParser()
	return &usecases{
		log:                    log,
		txm:                    txm,
		jwt:                    jwt,
		storage:                storage,
		uaparser:               uaparser,
		refreshExpiresDuration: refreshExpDuration,
	}
}
