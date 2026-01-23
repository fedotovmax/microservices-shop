package grpccontroller

import (
	"context"
	"errors"
	"log/slog"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/sessions_service/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Usecases interface {
	CreateSession(ctx context.Context, in *inputs.PrepareSessionInput) (*domain.SessionResponse, error)
	RefreshSession(ctx context.Context, in *inputs.RefreshSessionInput) (*domain.SessionResponse, error)
}

type controller struct {
	sessionspb.UnimplementedSessionsServiceServer
	log      *slog.Logger
	usecases Usecases
}

func New(log *slog.Logger, u Usecases) *controller {
	return &controller{
		log:      log,
		usecases: u,
	}
}

func (c *controller) handleError(locale string, fallback string, err error) error {

	const op = "controller.grpc.handleError"

	l := c.log.With(slog.String("op", op))

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
	default:
		l.Warn(err.Error())
		code = codes.Internal
		msgKey = fallback
	}

	msg, err := i18n.Local.Get(locale, msgKey)

	if err != nil {
		l.Warn("18n error", logger.Err(err))
	}

	st := status.New(code, msg)
	return st.Err()
}
