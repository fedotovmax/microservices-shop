package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
)

func (u *usecases) SendNewEmailVerifyLink(ctx context.Context, uid string, locale string) error {
	const op = "SendNewEmailVerifyLink"

	err := u.txm.Wrap(ctx, func(txctx context.Context) error {

		user, err := u.FindUserByID(txctx, uid)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if user.IsEmailVerified {
			return fmt.Errorf("%s: %w", op, errs.ErrUserEmailAlreadyVerified)
		}

		// link, err := u.FindEmailVerifyLinkByUserID(txctx, user.ID)

		// if err != nil {
		// 	return fmt.Errorf("%s: %w", op, err)
		// }

		expiresAt := time.Now().Add(u.cfg.EmailVerifyLinkExpiresDuration).UTC()

		newLink, err := u.emailVerifyStorage.UpdateByUserID(txctx, user.ID, expiresAt)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.createEmalVerifyLinkAddedEvent(txctx, &createEmalVerifyLinkAddedEventParams{
			ID:            user.ID,
			Email:         user.Email,
			Link:          newLink.Link,
			LinkExpiresAt: newLink.LinkExpiresAt,
			Locale:        locale,
		})

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil

	})

	return err
}
