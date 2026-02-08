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
	"github.com/fedotovmax/pgxtx"
	"github.com/google/uuid"
	"github.com/medama-io/go-useragent"
)

type createSessionData struct {
	uid            string
	browser        string
	browserVersion string
	os             string
	device         string
	ip             string
}

type CreateSessionUsecase struct {
	log                     *slog.Logger
	cfg                     *TokenConfig
	txm                     pgxtx.Manager
	uaParser                *useragent.Parser
	checkAllSecurityMethods *CheckAllSecurityMethodsUsecase
	sessionsStorage         ports.SessionsStorage
	securityStorage         ports.SecurityStorage
	usersQuery              queries.User
}

func NewCreateSessionUsecase(
	log *slog.Logger,
	cfg *TokenConfig,
	txm pgxtx.Manager,
	uaParser *useragent.Parser,
	checkAllSecurityMethods *CheckAllSecurityMethodsUsecase,
	sessionsStorage ports.SessionsStorage,
	securityStorage ports.SecurityStorage,
	usersQuery queries.User,
) *CreateSessionUsecase {
	return &CreateSessionUsecase{
		log:                     log,
		cfg:                     cfg,
		txm:                     txm,
		uaParser:                uaParser,
		checkAllSecurityMethods: checkAllSecurityMethods,
		sessionsStorage:         sessionsStorage,
		securityStorage:         securityStorage,
		usersQuery:              usersQuery,
	}
}

func (u *CreateSessionUsecase) Execute(pctx context.Context, in *inputs.PrepareSession) (*domain.SessionResponse, error) {

	const op = "usecases.security.CreateSession"

	agent := u.uaParser.Parse(in.GetUserAgent())

	if agent.IsBot() {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrAgentLooksLikeBot)
	}

	var newSession *domain.SessionResponse

	var independentTxErr error

	var basedTxErr error

	basedTxErr = u.txm.Wrap(pctx, func(txCtx context.Context) error {

		user, err := u.usersQuery.FindByID(txCtx, in.GetUID())

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if user.IsDeleted() {
			return fmt.Errorf("%s: %w", op, errs.ErrUserDeleted)
		}

		nowUTC := time.Now().UTC()

		preparedTrustToken, shouldRollback, err := u.checkAllSecurityMethods.Execute(
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

		_, err = u.sessionsStorage.Create(txCtx, &inputs.CreateSession{
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
				err = u.securityStorage.CreateTrustToken(txCtx, &inputs.CreateTrustToken{
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
				err = u.securityStorage.UpdateTrustToken(txCtx, &inputs.CreateTrustToken{
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
