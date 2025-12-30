package grpccontroller

import (
	"context"
	"errors"
	"log/slog"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Usecases interface {
	CreateSession(ctx context.Context, in *inputs.PrepareSessionInput) (*domain.SessionResponse, error)
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

	//TODO:
	switch {
	case errors.Is(err, errors.New("123")):
		code = 1
		msgKey = ""
	default:
		l.Warn(err.Error())
		code = codes.Internal
	}

	msg, err := i18n.Local.Get(locale, msgKey)

	if err != nil {
		l.Warn("18n error", logger.Err(err))
		msg = fallback
	}

	st := status.New(code, msg)
	return st.Err()
}
