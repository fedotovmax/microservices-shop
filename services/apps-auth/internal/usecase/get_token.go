package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/apps-auth/internal/adapter"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/utils"
	"github.com/fedotovmax/passport"
)

func (u *usecases) GetToken(ctx context.Context, secret string) (string, error) {
	const op = "usecases.GetToken"

	hash := utils.CreateHash(secret)

	app, err := u.storage.FindApp(ctx, hash)

	if err != nil {
		if errors.Is(err, adapter.ErrNotFound) {
			return "", fmt.Errorf("%s: %w: %v", op, errs.ErrAppNotFound, err)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	uid := fmt.Sprintf("%s:%d", app.Name, app.Type)

	token, _, err := passport.CreateAccessToken(passport.CreateParms{
		Issuer:          u.cfg.Issuer,
		UID:             uid,
		SID:             hash,
		Secret:          u.cfg.TokenSecret,
		ExpiresDuration: u.cfg.TokenExpDuration,
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil

}
