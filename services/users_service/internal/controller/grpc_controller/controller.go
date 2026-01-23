package grpccontroller

import (
	"context"
	"errors"
	"log/slog"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Usecases interface {
	CreateUser(ctx context.Context, in *inputs.CreateUserInput, locale string) (string, error)
	UpdateUserProfile(ctx context.Context, in *inputs.UpdateUserInput, locale string) error
	FindUserByID(ctx context.Context, id string) (*domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UserSessionAction(ctx context.Context, in *inputs.SessionActionInput) (*domain.UserOKResponse, error)
}

type controller struct {
	userspb.UnimplementedUserServiceServer
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
	case errors.Is(err, errs.ErrUserNotFound):
		code = codes.NotFound
		msgKey = keys.UserNotFound
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
		msg = fallback
	}

	st := status.New(code, msg)
	return st.Err()
}
