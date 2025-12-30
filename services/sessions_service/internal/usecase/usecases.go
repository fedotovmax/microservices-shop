package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log/slog"
	"math/big"
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

	AddToBlackList(ctx context.Context, in *inputs.SecurityInput) error
	UpdateBlacklistCode(ctx context.Context, in *inputs.SecurityInput) error
	RemoveUserFromBlacklist(ctx context.Context, uid string) error

	FindUser(ctx context.Context, uid string) (*domain.SessionsUser, error)
	FindUserSessions(ctx context.Context, uid string) ([]*domain.Session, error)

	AddIPBypass(ctx context.Context, in *inputs.SecurityInput) error
	UpdateIPBypass(ctx context.Context, in *inputs.SecurityInput) error

	SetEventStatusDone(ctx context.Context, id string) error
	SetEventsReservedToByIDs(ctx context.Context, ids []string, dur time.Duration) error
	RemoveEventReserve(ctx context.Context, id string) error
	CreateEvent(ctx context.Context, d *inputs.CreateEvent) (string, error)
	FindNewAndNotReservedEvents(ctx context.Context, limit int) ([]*domain.Event, error)
}

//todo: implement ip-bypass

type Config struct {
	RefreshExpiresDuration   time.Duration
	BlacklistCodeLength      uint8
	BlacklistCodeExpDuration time.Duration
	LoginBypassCodeLength    uint8
	LoginBypassExpDuration   time.Duration
}

type usecases struct {
	log      *slog.Logger
	txm      pgxtx.Manager
	jwt      JwtAdapter
	storage  Storage
	uaparser *useragent.Parser
	cfg      *Config
}

func New(log *slog.Logger, txm pgxtx.Manager, jwt JwtAdapter, storage Storage, cfg *Config) *usecases {
	uaparser := useragent.NewParser()
	return &usecases{
		log:      log,
		txm:      txm,
		jwt:      jwt,
		storage:  storage,
		uaparser: uaparser,
		cfg:      cfg,
	}
}

type createRefreshTokenResult struct {
	nohashed string
	hashed   string
}

func (u *usecases) createRefreshToken() (*createRefreshTokenResult, error) {
	refreshToken, err := u.generateRefreshToken()

	if err != nil {
		return nil, err
	}

	refreshHash := u.hashToken(refreshToken)

	resulst := &createRefreshTokenResult{
		nohashed: refreshToken,
		hashed:   refreshHash,
	}

	return resulst, nil
}

func (u *usecases) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (u *usecases) hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (u *usecases) generateSecurityCode(length uint8) (string, error) {

	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}

	return string(b), nil
}
