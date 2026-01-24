package security

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/utils"
)

type handleSecurityMethodsParams struct {
	User       *domain.SessionsUser
	TrustToken string
	BypassCode string
	NowUTC     time.Time
}

// first returned arg is struct *domain.PreparedTrustToken, can be nil
// second is should rollback transaction flag, needed for save side-effects after transaction
// third is err error
func (u *usecases) handleSecurityMethods(ctx context.Context, params handleSecurityMethodsParams) (*domain.PreparedTrustToken, bool, error) {

	const op = "usecases.security.handleSecurityMethods"

	if params.User.IsDeleted() {
		return nil, true, fmt.Errorf("%s: %w", op, errs.ErrUserDeleted)
	}

	err := u.handleUserBlacklist(ctx, params.User)

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

			trustTokenFromDB, err = u.findTrustToken(ctx, params.User.Info.UID, hashedTrustToken)

			isTrustTokenNotFound = errors.Is(err, errs.ErrTrustTokenNotFound)

			if err != nil && !isTrustTokenNotFound {
				return nil, true, fmt.Errorf("%s: %w", op, err)
			}
		}

		if isTrustTokenNotFound {
			if params.User.HasBypass() {
				err = u.checkActiveBypass(ctx, params.User, params.BypassCode)
				if err != nil {
					if errors.Is(err, errs.ErrBypassCodeExpired) {
						return nil, false, fmt.Errorf("%s: %w", op, err)
					}
					return nil, true, fmt.Errorf("%s: %w", op, err)
				}
			} else {
				var codeExpiresAt *time.Time
				codeExpiresAt, err = u.AddLoginIPBypass(ctx, params.User)
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

		if params.User.HasBypass() && params.BypassCode != "" {

			err = u.checkActiveBypass(ctx, params.User, params.BypassCode)

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

/*

		isTrustTokenNotFound := true

		var trustTokenFromDB *domain.DeviceTrustToken

		if params.TrustToken != "" {
			hashedTrustToken := utils.HashToken(params.TrustToken)

			trustTokenFromDB, err = u.findTrustToken(ctx, params.User.Info.UID, hashedTrustToken)

			isTrustTokenNotFound = errors.Is(err, errs.ErrTrustTokenNotFound)

			if err != nil && !isTrustTokenNotFound {
				return nil, true, fmt.Errorf("%s: %w", op, err)
			}
		}

			if isTrustTokenNotFound {
				....
			}


=============================

	hashedTrustToken := utils.HashToken(params.TrustToken)

	trustTokenFromDB, err := u.findTrustToken(ctx, params.User.Info.UID, hashedTrustToken)

	isTrustTokenNotFound := errors.Is(err, errs.ErrTrustTokenNotFound)

	if err != nil && !isTrustTokenNotFound {
		return nil, true, fmt.Errorf("%s: %w", op, err)
	}

				if isTrustTokenNotFound {
				....
			}

*/
