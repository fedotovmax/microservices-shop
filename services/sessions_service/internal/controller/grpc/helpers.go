package grpc

import (
	"errors"
	"log/slog"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/sessions_service/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func handleError(log *slog.Logger, locale string, fallback string, err error) error {

	var (
		code   codes.Code
		msgKey string
	)

	switch {

	case errors.Is(err, errs.ErrSessionNotFound):
		code = codes.NotFound
		msgKey = keys.SessionNotFound
	case errors.Is(err, errs.ErrSessionExpired):
		code = codes.Unauthenticated
		msgKey = keys.InvalidTokenOrExpired
	case errors.Is(err, errs.ErrUserNotFound):
		code = codes.NotFound
		msgKey = keys.UserNotFound
	case errors.Is(err, errs.ErrUserDeleted):
		code = codes.PermissionDenied
		msgKey = keys.UserDeleted
	case errors.Is(err, errs.ErrUserAlreadyExists):
		code = codes.AlreadyExists
		msgKey = keys.UserAlreadyExists
	case errors.Is(err, errs.ErrAgentLooksLikeBot):
		code = codes.PermissionDenied
		msgKey = keys.UserAgentLooksLikeBot
	default:
		log.Warn(err.Error())
		code = codes.Internal
		msgKey = fallback
	}

	msg, err := i18n.Local.Get(locale, msgKey)

	if err != nil {
		log.Warn("18n error", logger.Err(err))
	}

	st := status.New(code, msg)
	return st.Err()
}

func handleCreateSessionError(log *slog.Logger, locale string, fallback string, err error) (*sessionspb.CreateSessionResponse, error) {

	var newDeviceLoginErr *errs.LoginFromNewIPOrDeviceError
	var userInBlacklistErr *errs.UserSessionsInBlacklistError

	switch {

	case errors.Is(err, errs.ErrBadBypassCode):

		msg, i18nerr := i18n.Local.Get(locale, keys.InvalidCode)

		if i18nerr != nil {
			log.Warn("18n error", logger.Err(err))
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
			log.Warn("18n error", logger.Err(err))
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
			log.Warn("18n error", logger.Err(err))
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
		return nil, handleError(log, locale, fallback, err)
	}
}
