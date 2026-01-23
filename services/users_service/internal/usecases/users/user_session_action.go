package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
)

func (u *usecases) UserSessionAction(ctx context.Context, in *inputs.SessionActionInput) (
	*domain.UserSessionActionResponse, error) {

	const op = "usecases.users.UserSessionAction"

	user, err := u.FindUserByEmail(ctx, in.GetEmail())

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return domain.NewUserSessionActionResponse("", "", domain.UserSessionStatusBadCredentials), nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ok := comparePassword(in.GetPassword(), user.PasswordHash)

	if !ok {
		return domain.NewUserSessionActionResponse("", "", domain.UserSessionStatusBadCredentials), nil
	}

	if !user.IsEmailVerified {
		return domain.NewUserSessionActionResponse(user.ID, user.Email, domain.UserSessionStatusEmailNotVerified), nil
	}

	if user.DeletedAt != nil {
		return domain.NewUserSessionActionResponse(user.ID, user.Email, domain.UserSessionStatusDeleted), nil
	}

	return domain.NewUserSessionActionResponse(user.ID, user.Email, domain.UserSessionStatusOK), nil
}
