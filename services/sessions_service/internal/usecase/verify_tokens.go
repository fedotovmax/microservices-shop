package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/passport"
)

func (u *usecases) VerifyRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {

	const op = "usecases.VerifyRefreshToken"

	refreshTokenHash := u.hashToken(refreshToken)

	session, err := u.FindSessionByHash(ctx, refreshTokenHash)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if session.IsExpired() {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrSessionExpired)
	}

	err = u.handleUserBlacklist(ctx, session.User)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = u.handleSessionRevoked(ctx, session)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil

}

func (u *usecases) RefreshTokens(ctx context.Context, in *inputs.RefreshSessionInput) (*domain.SessionResponse, error) {

	const op = "usecases.RefreshTokens"

	agent := u.uaparser.Parse(in.GetUserAgent())

	if agent.IsBot() {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrAgentLooksLikeBot)
	}

	session, err := u.VerifyRefreshToken(ctx, in.GetRefreshToken())

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

	newRefreshToken, err := u.createRefreshToken()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	refreshExpTime := time.Now().Add(u.cfg.RefreshExpiresDuration).UTC()

	err = u.storage.sessions.UpdateSession(ctx, &inputs.CreateSessionInput{
		SID:            session.ID,
		UID:            session.User.Info.UID,
		RefreshHash:    newRefreshToken.hashed,
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
		RefreshToken:   newRefreshToken.nohashed,
		RefreshExpTime: refreshExpTime,
	}

	return response, nil
}
