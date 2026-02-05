package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/utils"
	"github.com/fedotovmax/passport"
)

func (u *usecases) RefreshSession(ctx context.Context, in *inputs.RefreshSession) (*domain.SessionResponse, error) {

	const op = "usecases.security.RefreshSession"

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
