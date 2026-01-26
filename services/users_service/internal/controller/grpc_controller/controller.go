package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
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
