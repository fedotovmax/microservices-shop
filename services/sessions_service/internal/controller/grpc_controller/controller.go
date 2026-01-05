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
	CreateSession(ctx context.Context, in *inputs.PrepareSessionInput, bypassCode string) (*domain.SessionResponse, error)
	RefreshTokens(ctx context.Context, in *inputs.RefreshSessionInput) (*domain.SessionResponse, error)
	VerifyAccessToken(ctx context.Context, in *inputs.VerifyAccessInput) (*domain.Session, error)
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

func handleError(l *slog.Logger, locale string, fallback string, err error) error {
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
	case errors.Is(err, errs.ErrLoginFromNewIPOrDevice):
		code = codes.PermissionDenied
		msgKey = keys.LoginFromNewIPOrDevice
	case errors.Is(err, errs.ErrBadBypassCode):
		code = codes.PermissionDenied
		msgKey = keys.BadBypassCode
	case errors.Is(err, errs.ErrBadBlacklistCode):
		code = codes.PermissionDenied
		msgKey = keys.BadBlacklistCode
	//TODO: add user in blacklist, user bypass, codes expired, bad codes, login from new ip or device
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
