package security

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/utils"
	"github.com/fedotovmax/passport"

	"github.com/google/uuid"
)

type createSessionData struct {
	uid            string
	browser        string
	browserVersion string
	os             string
	device         string
	ip             string
}

func (u *usecases) CreateSession(pctx context.Context, in *inputs.PrepareSessionInput) (*domain.SessionResponse, error) {

	const op = "usecases.security.CreateSession"

	agent := u.uaparser.Parse(in.GetUserAgent())

	if agent.IsBot() {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrAgentLooksLikeBot)
	}

	var newSession *domain.SessionResponse

	var independentTxErr error

	var basedTxErr error

	basedTxErr = u.txm.Wrap(pctx, func(txCtx context.Context) error {

		user, err := u.FindUserByID(txCtx, in.GetUID())

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if user.IsDeleted() {
			return fmt.Errorf("%s: %w", op, errs.ErrUserDeleted)
		}

		nowUTC := time.Now().UTC()

		preparedTrustToken, shouldRollback, err := u.handleSecurityMethods(
			txCtx,
			handleSecurityMethodsParams{
				User:       user,
				TrustToken: in.GetDeviceTrustToken(),
				BypassCode: in.GetBypassCode(),
				NowUTC:     nowUTC,
			},
		)

		if err != nil {
			if shouldRollback {
				return err
			}
			independentTxErr = err
			return nil
		}

		data := &createSessionData{
			uid:            user.Info.UID,
			browser:        agent.Browser().String(),
			browserVersion: agent.BrowserVersion(),
			os:             agent.OS().String(),
			device:         agent.Device().String(),
			ip:             in.GetIP(),
		}

		sid := uuid.New().String()

		refreshToken, err := utils.CreateToken()

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		refreshExpTime := nowUTC.Add(u.cfg.RefreshExpiresDuration)

		token, exp, err := passport.CreateAccessToken(passport.CreateParms{
			Issuer:          u.cfg.TokenIssuer,
			Secret:          u.cfg.TokenSecret,
			ExpiresDuration: u.cfg.AccessExpiresDuration,
			UID:             data.uid,
			SID:             sid,
		})

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		_, err = u.storage.CreateSession(txCtx, &inputs.CreateSessionInput{
			SID:            sid,
			UID:            data.uid,
			RefreshHash:    refreshToken.Hashed,
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

		var responseTrustToken *domain.SessionResponseTrustToken

		if preparedTrustToken != nil {
			switch preparedTrustToken.Action {
			case domain.TrustTokenCreated:
				err = u.storage.CreateTrustToken(txCtx, &inputs.CreateTrustTokenInput{
					TokenHash: preparedTrustToken.DeviceTrustTokenHash,
					UID:       user.Info.UID,
					ExpiresAt: preparedTrustToken.DeviceTrustTokenExpTime,
				})

				if err != nil {
					return fmt.Errorf("%s: %w", op, err)
				}
				responseTrustToken = &domain.SessionResponseTrustToken{
					DeviceTrustTokenExpTime: preparedTrustToken.DeviceTrustTokenExpTime,
					DeviceTrustTokenValue:   preparedTrustToken.DeviceTrustTokenValue,
				}
			case domain.TrustTokenUpdated:
				err = u.storage.UpdateTrustToken(txCtx, &inputs.CreateTrustTokenInput{
					TokenHash: preparedTrustToken.DeviceTrustTokenHash,
					UID:       preparedTrustToken.UID,
					ExpiresAt: preparedTrustToken.DeviceTrustTokenExpTime,
				})
				if err != nil {
					return fmt.Errorf("%s: %w", op, err)
				}
				responseTrustToken = &domain.SessionResponseTrustToken{
					DeviceTrustTokenExpTime: preparedTrustToken.DeviceTrustTokenExpTime,
					DeviceTrustTokenValue:   preparedTrustToken.DeviceTrustTokenValue,
				}
			}
		}

		newSession = &domain.SessionResponse{
			AccessToken:    token,
			RefreshToken:   refreshToken.Nohashed,
			AccessExpTime:  exp,
			RefreshExpTime: refreshExpTime,
			TrustToken:     responseTrustToken,
		}

		return nil
	})

	if basedTxErr != nil {
		return nil, basedTxErr
	}

	if independentTxErr != nil {
		return nil, independentTxErr
	}

	return newSession, nil
}
