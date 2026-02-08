package grpccontroller

import (
	"errors"
	"log/slog"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleError(log *slog.Logger, locale string, fallback string, err error) error {

	var (
		code   codes.Code
		msgKey string
	)

	switch {
	case errors.Is(err, errs.ErrUserNotFound):
		code = codes.NotFound
		msgKey = keys.UserNotFound
	case errors.Is(err, errs.ErrUserAlreadyExists):
		code = codes.AlreadyExists
		msgKey = keys.UserAlreadyExists
	case errors.Is(err, errs.ErrUserEmailAlreadyVerified):
		code = codes.AlreadyExists
		msgKey = keys.UserEmailAlreadyVerified
	default:
		log.Warn(err.Error())
		code = codes.Internal
		msgKey = fallback
	}

	msg, err := i18n.Local.Get(locale, msgKey)

	if err != nil {
		log.Warn("18n error", logger.Err(err))
		msg = fallback
	}

	st := status.New(code, msg)
	return st.Err()
}
