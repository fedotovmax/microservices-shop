package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
)

func (u *usecases) VerifyAccessToken(ctx context.Context, in *inputs.VerifyAccessInput) (*domain.Session, error) {

	const op = "usecases.VerifyAccessToken"

	sid, uid, err := u.jwt.ParseAccessToken(in.GetAccessToken(), in.GetIssuer())

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//TODO: switch errors?
	session, err := u.FindSessionByID(ctx, sid)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if session.User.Info.UID != uid {
		u.log.Warn("the received user ID is not equal to the session user ID")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrSessionNotFound)
	}

	return session, nil
}

func (u *usecases) VerifyRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {

	const op = "usecases.VerifyRefreshToken"

	refreshTokenHash := u.hashToken(refreshToken)

	session, err := u.FindSessionByHash(ctx, refreshTokenHash)

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

	//TODO: switch errors?
	session, err := u.VerifyRefreshToken(ctx, in.GetRefreshToken())

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	newAccessToken, err := u.jwt.CreateAccessToken(
		in.GetIssuer(),
		session.User.Info.UID,
		session.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	newRefreshToken, err := u.createRefreshToken()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	refreshExpTime := time.Now().Add(u.cfg.RefreshExpiresDuration)

	err = u.storage.UpdateSession(ctx, &inputs.CreateSessionInput{
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
		AccessToken:    newAccessToken.AccessToken,
		AccessExpTime:  newAccessToken.AccessExpTime,
		RefreshToken:   newRefreshToken.nohashed,
		RefreshExpTime: refreshExpTime,
	}

	return response, nil
}
