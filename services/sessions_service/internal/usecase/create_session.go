package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/google/uuid"
)

type createSessionData struct {
	issuer         string
	uid            string
	browser        string
	browserVersion string
	os             string
	device         string
	ip             string
}

func (u *usecases) CreateSession(pctx context.Context, in *inputs.PrepareSessionInput, bypassCode string) (*domain.SessionResponse, error) {

	const op = "usecases.CreateSession"

	agent := u.uaparser.Parse(in.GetUserAgent())

	if agent.IsBot() {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrAgentLooksLikeBot)
	}

	var newSession *domain.SessionResponse
	var err error

	err = u.txm.Wrap(pctx, func(txCtx context.Context) error {

		user, err := u.FindUserByID(txCtx, in.GetUID())

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.handleUserBlacklist(txCtx, user)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.handleUserBypass(txCtx, user, bypassParams{
			IP:         in.GetIP(),
			Browser:    agent.Browser().String(),
			BypassCode: bypassCode,
		})

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		data := &createSessionData{
			issuer:         in.GetIssuer(),
			uid:            user.Info.UID,
			browser:        agent.Browser().String(),
			browserVersion: agent.BrowserVersion(),
			os:             agent.OS().String(),
			device:         agent.Device().String(),
			ip:             in.GetIP(),
		}

		sid := uuid.New().String()

		refreshToken, err := u.createRefreshToken()

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		refreshExpTime := time.Now().Add(u.cfg.RefreshExpiresDuration)

		newAccessToken, err := u.jwt.CreateAccessToken(
			data.issuer,
			data.uid,
			sid,
		)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		_, err = u.storage.CreateSession(txCtx, &inputs.CreateSessionInput{
			SID:            sid,
			UID:            data.uid,
			RefreshHash:    refreshToken.hashed,
			Browser:        data.browser,
			BrowserVersion: data.browserVersion,
			OS:             data.os,
			Device:         data.device,
			IP:             data.ip,
			ExpiresAt:      refreshExpTime,
		})

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		newSession = &domain.SessionResponse{
			AccessToken:    newAccessToken.AccessToken,
			RefreshToken:   refreshToken.nohashed,
			AccessExpTime:  newAccessToken.AccessExpTime,
			RefreshExpTime: refreshExpTime,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return newSession, nil
}
