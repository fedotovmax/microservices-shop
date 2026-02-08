package grpccontroller

import (
	"log/slog"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/usecases"
)

type session struct {
	sessionspb.UnimplementedSessionsServiceServer
	log            *slog.Logger
	createSession  *usecases.CreateSessionUsecase
	refreshSession *usecases.RefreshSessionUsecase
}

func NewSession(
	log *slog.Logger,
	createSession *usecases.CreateSessionUsecase,
	refreshSession *usecases.RefreshSessionUsecase,
) *session {
	return &session{
		log:            log,
		createSession:  createSession,
		refreshSession: refreshSession,
	}
}
