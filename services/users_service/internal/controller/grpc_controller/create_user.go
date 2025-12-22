package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/pkg/utils/grpchelper"
)

func (c *grpcController) CreateUser(ctx context.Context, req *userspb.CreateUserRequest) (*userspb.CreateUserResponse, error) {

	const op = "controller.grpc.CreateUser"

	l := c.log.With(slog.String("op", op))
	metaParams := inputs.NewMetaParams()

	metaParams.SetLocale(grpchelper.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0])

	createUserInput := inputs.NewCreateUserInput()
	createUserInput.SetEmail(req.GetEmail())
	createUserInput.SetPassword(req.GetPassword())

	err := createUserInput.Validate(metaParams.GetLocale())

	if err != nil {
		return nil, grpchelper.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	userId, err := c.usecases.CreateUser(ctx, metaParams, createUserInput)

	if err != nil {
		return nil, handleError(l, metaParams.GetLocale(), keys.CreateUserInternal, err)
	}

	return &userspb.CreateUserResponse{
		Id: userId,
	}, nil
}
