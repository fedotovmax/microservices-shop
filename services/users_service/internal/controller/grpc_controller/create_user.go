package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/grpcutils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
)

func (c *controller) CreateUser(ctx context.Context, req *userspb.CreateUserRequest) (*userspb.CreateUserResponse, error) {

	const op = "controller.grpc.CreateUser"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	createUserInput := inputs.NewCreateUser()
	createUserInput.SetEmail(req.GetEmail())
	createUserInput.SetPassword(req.GetPassword())

	err := createUserInput.Validate(locale)

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	userId, err := c.usecases.CreateUser(ctx, createUserInput, locale)

	if err != nil {
		return nil, c.handleError(locale, keys.CreateUserInternal, err)
	}

	return &userspb.CreateUserResponse{
		Id: userId,
	}, nil
}
