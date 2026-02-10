package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/internal/ports"
	"github.com/fedotovmax/microservices-shop/users_service/internal/queries"
	"github.com/fedotovmax/microservices-shop/users_service/internal/utils"
	"github.com/fedotovmax/pgxtx"
)

type SessionActionUsecase struct {
	txm pgxtx.Manager
	log *slog.Logger

	usersStorage ports.UsersStorage
	query        queries.Users
}

func NewSessionActionUsecase(
	txm pgxtx.Manager,
	log *slog.Logger,
	usersStorage ports.UsersStorage,
	query queries.Users,
) *SessionActionUsecase {
	return &SessionActionUsecase{
		txm:          txm,
		log:          log,
		usersStorage: usersStorage,
		query:        query,
	}
}

func (u *SessionActionUsecase) Execute(ctx context.Context, in *inputs.SessionAction) (
	*domain.UserOKResponse, error) {

	const op = "usecases.session_action"

	user, err := u.query.FindByEmail(ctx, in.GetEmail())

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, fmt.Errorf("%s: %w: %v", op, errs.ErrBadCredentials, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ok := utils.ComparePasswords(in.GetPassword(), user.PasswordHash)

	if !ok {
		return nil, fmt.Errorf("%s: %w: %v", op, errs.ErrBadCredentials, err)
	}

	if user.DeletedAt != nil {
		//todo: change last chance to restore
		deletedErr := errs.NewUserDeletedError(keys.UserDeleted, *user.DeletedAt, time.Now().UTC().Add(time.Hour*730))
		return nil, fmt.Errorf("%s: %w", op, deletedErr)
	}

	if !user.IsEmailVerified {
		return nil, fmt.Errorf("%s: %w: %v", op, errs.NewEmailNotVerifiedErrorError(
			keys.UserEmailNotVerified,
			user.ID,
		), err)
	}

	return &domain.UserOKResponse{
		UID:   user.ID,
		Email: user.Email,
	}, nil
}
