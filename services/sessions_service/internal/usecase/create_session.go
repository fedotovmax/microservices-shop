package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/google/uuid"
)

func (u *usecases) CreateSession(ctx context.Context, in *inputs.PrepareSessionInput) (*domain.SessionResponse, error) {

	const op = "usecases.CreateSession"

	agent := u.uaparser.Parse(in.GetUserAgent())

	if agent.IsBot() {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrAgentLooksLikeBot)
	}

	user, err := u.storage.FindUser(ctx, in.GetUID())

	if err != nil && !errors.Is(err, adapter.ErrNotFound) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err == nil {
		return nil, fmt.Errorf("%s: %w", op, errs.NewUserSessionsInBlacklistError(
			user.Info.Email, user.Info.UID,
		))
	}

	userActiveSessions, err := u.storage.FindUserSessions(ctx, in.GetUID())
	_, _ = userActiveSessions, err

	sid := uuid.New().String()

	refreshToken, err := u.createRefreshToken()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	refreshExpTime := time.Now().Add(u.refreshExpiresDuration)

	newAccessToken, err := u.jwt.CreateAccessToken(
		in.GetIssuer(),
		in.GetUID(),
		sid,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = u.storage.CreateSession(ctx, &inputs.CreateSessionInput{
		SID:            sid,
		UID:            in.GetUID(),
		RefreshHash:    refreshToken.hashed,
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

	return &domain.SessionResponse{
		AccessToken:    newAccessToken.AccessToken,
		RefreshToken:   refreshToken.nohashed,
		AccessExpTime:  newAccessToken.AccessExpTime,
		RefreshExpTime: refreshExpTime,
	}, nil
}
