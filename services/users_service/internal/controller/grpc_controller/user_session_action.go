package grpccontroller

import (
	"context"
	"errors"
	"log/slog"

	"github.com/fedotovmax/goutils/timeutils"
	"github.com/fedotovmax/grpcutils"
	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/pkg/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *controller) UserSessionAction(ctx context.Context, req *userspb.UserSessionActionRequest) (*userspb.UserSessionActionResponse, error) {

	const op = "controller.grpc.UserSessionAction"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	userSessionActionInput := inputs.NewSessionActionInput()
	userSessionActionInput.SetEmail(req.GetEmail())
	userSessionActionInput.SetPassword(req.GetPassword())

	err := userSessionActionInput.Validate(locale)

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	sessionActionResponse, err := c.usecases.UserSessionAction(ctx, userSessionActionInput)

	if err != nil {
		return c.handleSessionActionError(locale, keys.UserSessionActionInternal, err)
	}

	return &userspb.UserSessionActionResponse{
		Payload: &userspb.UserSessionActionResponse_Ok{
			Ok: &userspb.UserOK{
				Email:  sessionActionResponse.Email,
				UserId: sessionActionResponse.UID,
			},
		},
	}, nil
}

func (c *controller) handleSessionActionError(locale string, fallbackMsg string, err error) (*userspb.UserSessionActionResponse, error) {

	const op = "controller.grpc.handleSessionActionError"

	l := c.log.With(slog.String("op", op))

	var deletedErr *errs.UserDeletedError
	var emailNotVerifiedErr *errs.EmailNotVerifiedError

	switch {

	case errors.As(err, &deletedErr):

		var formattedTime string

		if locale == keys.RuLocale {
			formattedTime = timeutils.FormatDateRU(deletedErr.LastChanceRestore)
		} else {
			formattedTime = timeutils.TimeToString(keys.EnShortDateFormat, deletedErr.LastChanceRestore)
		}

		msg, i18nerr := i18n.Local.Get(locale, deletedErr.ErrCode, formattedTime)

		if i18nerr != nil {
			l.Warn("18n error", logger.Err(err))
		}

		return &userspb.UserSessionActionResponse{
			Payload: &userspb.UserSessionActionResponse_Deleted{
				Deleted: &userspb.UserDeleted{
					Message:           msg,
					DeletedAt:         timestamppb.New(deletedErr.DeletedAt),
					LastChanceRestore: timestamppb.New(deletedErr.LastChanceRestore),
				},
			},
		}, nil

	case errors.Is(err, errs.ErrBadCredentials):

		msg, i18nerr := i18n.Local.Get(locale, keys.UserBadCredentials)

		if i18nerr != nil {
			l.Warn("18n error", logger.Err(err))
		}

		return &userspb.UserSessionActionResponse{
			Payload: &userspb.UserSessionActionResponse_BadCredentials{
				BadCredentials: &userspb.UserBadCredentials{
					Message: msg,
				},
			},
		}, nil

	case errors.As(err, &emailNotVerifiedErr):
		msg, i18nerr := i18n.Local.Get(locale, emailNotVerifiedErr.ErrCode)

		if i18nerr != nil {
			l.Warn("18n error", logger.Err(err))
		}

		return &userspb.UserSessionActionResponse{
			Payload: &userspb.UserSessionActionResponse_EmailNotVerified{
				EmailNotVerified: &userspb.UserEmailNotVerified{
					Message: msg,
					UserId:  emailNotVerifiedErr.UID,
				},
			},
		}, nil

	default:
		return nil, c.handleError(locale, fallbackMsg, err)
	}
}
