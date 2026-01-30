package users

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
)

func (u *usecases) VerifyEmail(ctx context.Context, link string) error {

	const op = "usecase.users.VerifyEmail"

	err := u.txm.Wrap(ctx, func(txctx context.Context) error {

		linkEntity, err := u.FindEmailVerifyLinkByPrimary(txctx, link)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if linkEntity.IsExpired() {
			return fmt.Errorf("%s: %w", op, errs.NewVerifyEmailLinkExpiredError(
				keys.VerifyEmailLinkExpired,
				linkEntity.UserID,
			))
		}

		err = u.storage.SetIsEmailVerified(txctx, linkEntity.UserID, true)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.storage.DeleteEmailVerifyLink(txctx, linkEntity.Link)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil

	})

	return err
}
