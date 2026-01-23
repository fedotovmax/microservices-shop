package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/grpcutils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
)

func (c *controller) FindUserByID(ctx context.Context, req *userspb.FindUserByIDRequest) (*userspb.User, error) {

	const op = "controller.grpc.FindUserByID"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	findUserByIdInput := inputs.NewFindUserByIDInput()
	findUserByIdInput.SetUserID(req.GetId())

	err := findUserByIdInput.Validate(locale)

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	user, err := c.usecases.FindUserByID(ctx, req.GetId())

	if err != nil {
		return nil, c.handleError(locale, keys.GetUserInternal, err)
	}

	return user.ToProto(), nil

}
