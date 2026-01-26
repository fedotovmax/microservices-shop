package grpccontroller

import (
	"errors"
	"log/slog"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/sessions_service/pkg/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *controller) handleCreateSessionError(locale string, fallback string, err error) (*sessionspb.CreateSessionResponse, error) {
	const op = "controller.grpc.handleCreateSessionError"

	l := c.log.With(slog.String("op", op))

	var newDeviceLoginErr *errs.LoginFromNewIPOrDeviceError
	var userInBlacklistErr *errs.UserSessionsInBlacklistError

	switch {

	case errors.Is(err, errs.ErrBadBypassCode):

		msg, i18nerr := i18n.Local.Get(locale, keys.InvalidCode)

		if i18nerr != nil {
			l.Warn("18n error", logger.Err(err))
		}

		return &sessionspb.CreateSessionResponse{
			Payload: &sessionspb.CreateSessionResponse_BadBypassCode{
				BadBypassCode: &sessionspb.BadBypassCode{
					Message: msg,
				},
			},
		}, nil

	case errors.As(err, &newDeviceLoginErr):
		msg, i18nerr := i18n.Local.Get(locale, newDeviceLoginErr.ErrCode)

		if i18nerr != nil {
			l.Warn("18n error", logger.Err(err))
		}

		return &sessionspb.CreateSessionResponse{
			Payload: &sessionspb.CreateSessionResponse_LoginFromNewDevice{
				LoginFromNewDevice: &sessionspb.LoginFromNewDeviceOrIP{
					Message:       msg,
					CodeExpiresAt: timestamppb.New(newDeviceLoginErr.CodeExpiresAt),
				},
			},
		}, nil

	case errors.As(err, &userInBlacklistErr):

		msg, i18nerr := i18n.Local.Get(locale, userInBlacklistErr.ErrCode)

		if i18nerr != nil {
			l.Warn("18n error", logger.Err(err))
		}

		return &sessionspb.CreateSessionResponse{
			Payload: &sessionspb.CreateSessionResponse_UserInBlacklist{
				UserInBlacklist: &sessionspb.UserInBlacklist{
					Message:       msg,
					LinkExpiresAt: timestamppb.New(userInBlacklistErr.LinkExpiresAt),
				},
			},
		}, nil

	default:
		return nil, c.handleError(locale, fallback, err)
	}
}
