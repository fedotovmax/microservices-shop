package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/appsauthpb"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/domain"
)

type Usecases interface {
	GetToken(ctx context.Context, secret string) (string, error)
	CreateApp(ctx context.Context, newApp *domain.App) (string, error)
}

type controller struct {
	appsauthpb.UnimplementedAppsAuthServiceServer
	log      *slog.Logger
	usecases Usecases
}

func New(log *slog.Logger, u Usecases) *controller {
	return &controller{
		log:      log,
		usecases: u,
	}
}
