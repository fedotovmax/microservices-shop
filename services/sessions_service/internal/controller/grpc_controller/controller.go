package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
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
