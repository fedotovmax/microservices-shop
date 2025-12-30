package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/goutils/sliceutils"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/adapter"
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

func (u *usecases) createSessionFn(ctx context.Context, data *createSessionData) (*domain.SessionResponse, error) {

	const op = "usecases.createSessionFn"

	sid := uuid.New().String()

	refreshToken, err := u.createRefreshToken()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	refreshExpTime := time.Now().Add(u.cfg.RefreshExpiresDuration)

	newAccessToken, err := u.jwt.CreateAccessToken(
		data.issuer,
		data.uid,
		sid,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = u.storage.CreateSession(ctx, &inputs.CreateSessionInput{
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
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.SessionResponse{
		AccessToken:    newAccessToken.AccessToken,
		RefreshToken:   refreshToken.nohashed,
		AccessExpTime:  newAccessToken.AccessExpTime,
		RefreshExpTime: refreshExpTime,
	}, nil
}

func (u *usecases) CreateSession(ctx context.Context, in *inputs.PrepareSessionInput) (*domain.SessionResponse, error) {

	//TODO: in tx?

	//TODO: get bypass code
	bypassCode := "123456789123"

	const op = "usecases.CreateSession"

	agent := u.uaparser.Parse(in.GetUserAgent())

	if agent.IsBot() {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrAgentLooksLikeBot)
	}

	user, err := u.storage.FindUser(ctx, in.GetUID())

	if err != nil && !errors.Is(err, adapter.ErrNotFound) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if user.IsInBlackList() {
		return nil, fmt.Errorf("%s: %w", op, errs.NewUserSessionsInBlacklistError(
			user.Info.Email, user.Info.UID,
		))
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

	if !user.HasBypass() {
		userActiveSessions, err := u.storage.FindUserSessions(ctx, in.GetUID())

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		_, isIPFound := sliceutils.Find(userActiveSessions, func(session *domain.Session) bool {
			return session.IP == in.GetIP()
		})

		if !isIPFound {
			return nil, fmt.Errorf("%s: %w", op, errs.NewLoginFromNewIPOrDeviceError(
				user.Info.Email, user.Info.UID,
			))
		}

		newSession, err := u.createSessionFn(ctx, data)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return newSession, nil
	}

	if user.Bypass.IsCodeExpired() {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrBypassCodeExpired)
	}

	if !user.Bypass.ComapreCodes(bypassCode) {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrBadBypassCode)
	}

	newSession, err := u.createSessionFn(ctx, data)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return newSession, nil

}
