package usecases

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	eventspublisher "github.com/fedotovmax/microservices-shop/users_service/internal/events_publisher"
	"github.com/fedotovmax/microservices-shop/users_service/internal/ports"
	"github.com/fedotovmax/microservices-shop/users_service/internal/queries"
	"github.com/fedotovmax/pgxtx"
)

type SendNewEmailVerifyLinkUsecase struct {
	txm               pgxtx.Manager
	log               *slog.Logger
	cfg               *EmailConfig
	usersStorage      ports.UsersStorage
	verifyLinkStorage ports.EmailVerifyStorage
	publisher         eventspublisher.Publisher
	query             queries.Users
}

func NewSendNewEmailVerifyLinkUsecase(
	txm pgxtx.Manager,
	log *slog.Logger,
	cfg *EmailConfig,
	usersStorage ports.UsersStorage,
	verifyLinkStorage ports.EmailVerifyStorage,
	publisher eventspublisher.Publisher,
	query queries.Users,
) *SendNewEmailVerifyLinkUsecase {
	return &SendNewEmailVerifyLinkUsecase{
		txm:               txm,
		log:               log,
		cfg:               cfg,
		usersStorage:      usersStorage,
		verifyLinkStorage: verifyLinkStorage,
		publisher:         publisher,
		query:             query,
	}
}

func (u *SendNewEmailVerifyLinkUsecase) Execute(
	ctx context.Context,
	uid string,
	locale string,
) error {

	const op = "usecases.send_new_email_verify_link"

	err := u.txm.Wrap(ctx, func(txctx context.Context) error {

		user, err := u.query.FindByID(txctx, uid)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if user.IsEmailVerified {
			return fmt.Errorf("%s: %w", op, errs.ErrUserEmailAlreadyVerified)
		}

		expiresAt := time.Now().Add(u.cfg.EmailVerifyLinkExpiresDuration).UTC()

		newLink, err := u.verifyLinkStorage.UpdateByUserID(txctx, user.ID, expiresAt)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.publisher.UserEmalVerifyLinkAdded(txctx, &eventspublisher.UserEmalVerifyLinkAddedParams{
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
