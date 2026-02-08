package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/queries"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/utils"
)

type handleSecurityMethodsParams struct {
	User       *domain.SessionsUser
	TrustToken string
	BypassCode string
	NowUTC     time.Time
}

type CheckAllSecurityMethodsUsecase struct {
	log               *slog.Logger
	cfg               *SecurityConfig
	isSessionRevoked  *IsSessionRevokedUsecase
	isUserInBlacklist *IsUserInBlacklistUsecase
	checkBypass       *CheckBypassUsecase
	addLoginBypass    *AddLoginBypassUsecase
	trustTokenQuery   queries.TrustToken
}

func NewCheckAllSecurityMethodsUsecase(
	log *slog.Logger,
	cfg *SecurityConfig,
	isSessionRevoked *IsSessionRevokedUsecase,
	isUserInBlacklist *IsUserInBlacklistUsecase,
	checkBypass *CheckBypassUsecase,
	addLoginBypass *AddLoginBypassUsecase,
	trustTokenQuery queries.TrustToken,
) *CheckAllSecurityMethodsUsecase {
	return &CheckAllSecurityMethodsUsecase{
		log:               log,
		cfg:               cfg,
		isSessionRevoked:  isSessionRevoked,
		isUserInBlacklist: isUserInBlacklist,
		checkBypass:       checkBypass,
		addLoginBypass:    addLoginBypass,
		trustTokenQuery:   trustTokenQuery,
	}
}

func (u *CheckAllSecurityMethodsUsecase) Execute(ctx context.Context, params handleSecurityMethodsParams) (*domain.PreparedTrustToken, bool, error) {

	const op = "usecases.check_all_security_methods"

	if params.User.IsDeleted() {
		return nil, true, fmt.Errorf("%s: %w", op, errs.ErrUserDeleted)
	}

	err := u.isUserInBlacklist.Execute(ctx, params.User)

	if err != nil {
		if errors.Is(err, errs.ErrUserSessionsInBlackList) || errors.Is(err, errs.ErrBlacklistCodeExpired) {
			return nil, false, fmt.Errorf("%s: %w", op, err)
		}
		return nil, true, fmt.Errorf("%s: %w", op, err)
	}

	var preparedTrustToken *domain.PreparedTrustToken = nil

	if !params.User.HasTwoFactor() {

		isTrustTokenNotFound := true

		var trustTokenFromDB *domain.DeviceTrustToken

		if params.TrustToken != "" {
			hashedTrustToken := utils.HashToken(params.TrustToken)

			trustTokenFromDB, err = u.trustTokenQuery.Find(ctx, params.User.Info.UID, hashedTrustToken)

			isTrustTokenNotFound = errors.Is(err, errs.ErrTrustTokenNotFound)

			if err != nil && !isTrustTokenNotFound {
				return nil, true, fmt.Errorf("%s: %w", op, err)
			}
		}

		if isTrustTokenNotFound {
			if params.User.HasBypass() {
				err = u.checkBypass.Execute(ctx, params.User, params.BypassCode)
				if err != nil {
					if errors.Is(err, errs.ErrBypassCodeExpired) {
						return nil, false, fmt.Errorf("%s: %w", op, err)
					}
					return nil, true, fmt.Errorf("%s: %w", op, err)
				}
			} else {
				var codeExpiresAt *time.Time
				codeExpiresAt, err = u.addLoginBypass.Execute(ctx, params.User)
				if err != nil {
					return nil, true, fmt.Errorf("%s: %w", op, err)
				}
				return nil, false, fmt.Errorf("%s: %w", op, errs.NewLoginFromNewIPOrDeviceError(
					keys.LoginFromNewIPOrDevice,
					*codeExpiresAt,
				))
			}

			newTrustToken, err := utils.CreateToken()

			if err != nil {
				return nil, true, fmt.Errorf("%s: %w", op, err)
			}

			newExpires := params.NowUTC.Add(u.cfg.DeviceTrustTokenExpDuration)

			preparedTrustToken = &domain.PreparedTrustToken{
				UID:                     params.User.Info.UID,
				DeviceTrustTokenValue:   newTrustToken.Nohashed,
				DeviceTrustTokenHash:    newTrustToken.Hashed,
				DeviceTrustTokenExpTime: newExpires,
				Action:                  domain.TrustTokenCreated,
			}

		} else {
			newExpires := utils.ExtendTrustTokenTTL(trustTokenFromDB.ExpiresAt, params.NowUTC, u.cfg.DeviceTrustTokenThreshold, u.cfg.DeviceTrustTokenExpDuration)

			preparedTrustToken = &domain.PreparedTrustToken{
				UID:                     params.User.Info.UID,
				DeviceTrustTokenValue:   params.TrustToken,
				DeviceTrustTokenHash:    trustTokenFromDB.TokenHash,
				DeviceTrustTokenExpTime: newExpires,
				Action:                  domain.TrustTokenUpdated,
			}
		}
	} else {

		//TODO: maybe not check bypass if has two factor??
		if params.User.HasBypass() && params.BypassCode != "" {
			err = u.checkBypass.Execute(ctx, params.User, params.BypassCode)

			if err != nil {
				if errors.Is(err, errs.ErrBypassCodeExpired) {
					return nil, false, fmt.Errorf("%s: %w", op, err)
				}
				return nil, true, fmt.Errorf("%s: %w", op, err)
			}
		} else {

			//TODO: check otp two factor)
		}
	}

	return preparedTrustToken, false, nil
}
