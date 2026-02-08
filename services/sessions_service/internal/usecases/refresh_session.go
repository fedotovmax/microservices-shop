package usecases

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/ports"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/queries"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/utils"
	"github.com/fedotovmax/passport"
	"github.com/medama-io/go-useragent"
)

type RefreshSessionUsecase struct {
	log               *slog.Logger
	cfg               *TokenConfig
	isUserInBlacklist *IsUserInBlacklistUsecase
	isSessionRevoked  *IsSessionRevokedUsecase
	uaParser          *useragent.Parser
	sessionsStorage   ports.SessionsStorage
	sessionsQuery     queries.Session
}

func NewRefreshSessionUsecase(
	log *slog.Logger,
	cfg *TokenConfig,
	isUserInBlacklist *IsUserInBlacklistUsecase,
	isSessionRevoked *IsSessionRevokedUsecase,
	uaParser *useragent.Parser,
	sessionsStorage ports.SessionsStorage,
	sessionsQuery queries.Session,
) *RefreshSessionUsecase {
	return &RefreshSessionUsecase{
		log:               log,
		cfg:               cfg,
		isUserInBlacklist: isUserInBlacklist,
		isSessionRevoked:  isSessionRevoked,
		uaParser:          uaParser,
		sessionsStorage:   sessionsStorage,
		sessionsQuery:     sessionsQuery,
	}
}

func (u *RefreshSessionUsecase) Execute(ctx context.Context, in *inputs.RefreshSession) (*domain.SessionResponse, error) {
	const op = "usecases.refresh_session"

	agent := u.uaParser.Parse(in.GetUserAgent())

	if agent.IsBot() {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrAgentLooksLikeBot)
	}

	session, err := u.verifyRefreshToken(ctx, in.GetRefreshToken())

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	token, exp, err := passport.CreateAccessToken(passport.CreateParms{
		Issuer:          u.cfg.TokenIssuer,
		Secret:          u.cfg.TokenSecret,
		ExpiresDuration: u.cfg.AccessExpiresDuration,
		UID:             session.User.Info.UID,
		SID:             session.ID,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	newRefreshToken, err := utils.CreateToken()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	refreshExpTime := time.Now().Add(u.cfg.RefreshExpiresDuration).UTC()

	err = u.sessionsStorage.Update(ctx, &inputs.CreateSession{
		SID:            session.ID,
		UID:            session.User.Info.UID,
		RefreshHash:    newRefreshToken.Hashed,
		Browser:        agent.Browser().String(),
		BrowserVersion: agent.BrowserVersion(),
		OS:             agent.OS().String(),
		Device:         agent.Device().String(),
		IP:             in.GetIP(),
		ExpiresAt:      refreshExpTime,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	response := &domain.SessionResponse{
		AccessToken:    token,
		AccessExpTime:  exp,
		RefreshToken:   newRefreshToken.Nohashed,
		RefreshExpTime: refreshExpTime,
	}

	return response, nil

}

func (u *RefreshSessionUsecase) verifyRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {

	refreshTokenHash := utils.HashToken(refreshToken)

	session, err := u.sessionsQuery.FindByHash(ctx, refreshTokenHash)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if session.User.IsDeleted() {
		return nil, fmt.Errorf("%w", errs.ErrUserDeleted)
	}

	if session.IsExpired() {
		return nil, fmt.Errorf("%w", errs.ErrSessionExpired)
	}

	err = u.isUserInBlacklist.Execute(ctx, session.User)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	err = u.isSessionRevoked.Execute(ctx, session)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return session, nil

}
