package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	eventspublisher "github.com/fedotovmax/microservices-shop/users_service/internal/events_publisher"
	"github.com/fedotovmax/microservices-shop/users_service/internal/ports"
	"github.com/fedotovmax/microservices-shop/users_service/internal/queries"
	"github.com/fedotovmax/pgxtx"
)

type UpdateProfileUsecase struct {
	txm          pgxtx.Manager
	log          *slog.Logger
	usersStorage ports.UsersStorage
	publisher    eventspublisher.Publisher
	query        queries.Users
}

func NewUpdateProfileUsecase(
	txm pgxtx.Manager,
	log *slog.Logger,
	usersStorage ports.UsersStorage,
	verifyLinkStorage ports.EmailVerifyStorage,
	publisher eventspublisher.Publisher,
	query queries.Users,
) *UpdateProfileUsecase {
	return &UpdateProfileUsecase{
		txm:          txm,
		log:          log,
		usersStorage: usersStorage,
		publisher:    publisher,
		query:        query,
	}
}

func (u *UpdateProfileUsecase) Execute(ctx context.Context, in *inputs.UpdateUser, locale string) error {

	const op = "usecase.update_profile"

	err := u.txm.Wrap(ctx, func(txCtx context.Context) error {

		user, err := u.query.FindByID(txCtx, in.GetUserID())

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.usersStorage.UpdateProfile(txCtx, user.ID, in)

		if err != nil && !errors.Is(err, adapters.ErrNoFieldsToUpdate) {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.publisher.ProfileUpdated(txCtx, &eventspublisher.ProfileUpdatedParams{
			UserID: user.ID,
			Email:  user.Email,
			Locale: locale,
			Input:  in,
		})

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	return err
}
