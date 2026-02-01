package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
)

func (u *usecases) UserSessionAction(ctx context.Context, in *inputs.SessionActionInput) (
	*domain.UserOKResponse, error) {

	const op = "usecases.users.UserSessionAction"

	user, err := u.FindUserByEmail(ctx, in.GetEmail())

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, fmt.Errorf("%s: %w: %v", op, errs.ErrBadCredentials, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ok := comparePassword(in.GetPassword(), user.PasswordHash)

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
