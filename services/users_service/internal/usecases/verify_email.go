package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/internal/ports"
	"github.com/fedotovmax/microservices-shop/users_service/internal/queries"
	"github.com/fedotovmax/pgxtx"
)

type VerifyEmailUsecase struct {
	txm               pgxtx.Manager
	log               *slog.Logger
	usersStorage      ports.UsersStorage
	verifyLinkStorage ports.EmailVerifyStorage
	query             queries.EmailVerifyLink
}

func NewVerifyEmailUsecase(
	txm pgxtx.Manager,
	log *slog.Logger,
	usersStorage ports.UsersStorage,
	verifyLinkStorage ports.EmailVerifyStorage,
	query queries.EmailVerifyLink,
) *VerifyEmailUsecase {
	return &VerifyEmailUsecase{
		txm:               txm,
		log:               log,
		usersStorage:      usersStorage,
		verifyLinkStorage: verifyLinkStorage,
		query:             query,
	}
}

func (u *VerifyEmailUsecase) Execute(ctx context.Context, link string) error {

	const op = "usecase.verify_email"

	err := u.txm.Wrap(ctx, func(txctx context.Context) error {

		linkEntity, err := u.query.Find(txctx, link)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if linkEntity.IsExpired() {
			return fmt.Errorf("%s: %w", op, errs.NewVerifyEmailLinkExpiredError(
				keys.VerifyEmailLinkExpired,
				linkEntity.UserID,
			))
		}

		err = u.usersStorage.SetIsEmailVerified(txctx, linkEntity.UserID, true)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.verifyLinkStorage.Delete(txctx, linkEntity.Link)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil

	})

	return err
}
