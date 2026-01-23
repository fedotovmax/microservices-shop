package security

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
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

		hashedTrustToken := utils.HashToken(params.TrustToken)

		trustTokenFromDB, err := u.findTrustToken(ctx, params.User.Info.UID, hashedTrustToken)

		isTrustTokenNotFound := errors.Is(err, errs.ErrTrustTokenNotFound)

		if err != nil && !isTrustTokenNotFound {
			return nil, true, fmt.Errorf("%s: %w", op, err)
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

				err = u.AddLoginIPBypass(ctx, params.User)
				if err != nil {
					return nil, true, fmt.Errorf("%s: %w", op, err)
				}
				return nil, false, fmt.Errorf("%s: %w", op, errs.ErrLoginFromNewIPOrDevice)
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


OLD create sessions security checks, do not delete!


if user.IsDeleted() {
			return fmt.Errorf("%s: %w", op, errs.ErrUserDeleted)
		}

		err = u.handleUserBlacklist(txCtx, user)

		if err != nil {
			if errors.Is(err, errs.ErrUserSessionsInBlackList) || errors.Is(err, errs.ErrBlacklistCodeExpired) {
				independentTxErr = fmt.Errorf("%s: %w", op, err)
				return nil
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		var preparedTrustToken *domain.PreparedTrustToken = nil

		nowUTC := time.Now().UTC()

		deviceTrustTokenFromInput := in.GetDeviceTrustToken()
		bypassCode := in.GetBypassCode()

		if !user.HasTwoFactor() {

			hashedTrustToken := u.hashToken(deviceTrustTokenFromInput)

			trustTokenFromDB, err := u.findTrustToken(txCtx, user.Info.UID, hashedTrustToken)

			isTrustTokenNotFound := errors.Is(err, errs.ErrTrustTokenNotFound)

			if err != nil && !isTrustTokenNotFound {
				return fmt.Errorf("%s: %w", op, err)
			}

			if isTrustTokenNotFound {

				if user.HasBypass() {
					err = u.checkActiveBypass(txCtx, user, bypassCode)
					if err != nil {
						if errors.Is(err, errs.ErrBypassCodeExpired) {
							independentTxErr = fmt.Errorf("%s: %w", op, err)
							return nil
						}
						return fmt.Errorf("%s: %w", op, err)
					}
				} else {
					err = u.AddLoginIPBypass(txCtx, user)
					if err != nil {
						return fmt.Errorf("%s: %w", op, err)
					}
					independentTxErr = fmt.Errorf("%s: %w", op, errs.ErrLoginFromNewIPOrDevice)
					return nil
				}

				newTrustToken, err := u.createToken()

				if err != nil {
					return fmt.Errorf("%s: %w", op, err)
				}

				newExpires := nowUTC.Add(u.cfg.DeviceTrustTokenExpDuration)

				preparedTrustToken = &domain.PreparedTrustToken{
					UID:                     user.Info.UID,
					DeviceTrustTokenValue:   newTrustToken.nohashed,
					DeviceTrustTokenHash:    newTrustToken.hashed,
					DeviceTrustTokenExpTime: newExpires,
					Action:                  domain.TrustTokenCreated,
				}

			} else {
				newExpires := u.extendTrustTokenTTL(trustTokenFromDB.ExpiresAt, nowUTC, u.cfg.DeviceTrustTokenThreshold, u.cfg.DeviceTrustTokenExpDuration)

				preparedTrustToken = &domain.PreparedTrustToken{
					UID:                     user.Info.UID,
					DeviceTrustTokenValue:   deviceTrustTokenFromInput,
					DeviceTrustTokenHash:    trustTokenFromDB.TokenHash,
					DeviceTrustTokenExpTime: newExpires,
					Action:                  domain.TrustTokenUpdated,
				}
			}
		} else {
			if user.HasBypass() && bypassCode != "" {

				err = u.checkActiveBypass(txCtx, user, bypassCode)

				if err != nil {
					if errors.Is(err, errs.ErrBypassCodeExpired) {
						independentTxErr = fmt.Errorf("%s: %w", op, err)
						return nil
					}
					return fmt.Errorf("%s: %w", op, err)
				}
			} else {
				//TODO: check otp two factor)
			}
		}

*/
