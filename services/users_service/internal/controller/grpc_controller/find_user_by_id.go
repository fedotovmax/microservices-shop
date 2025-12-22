package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/pkg/utils/grpchelper"
)

func (c *grpcController) FindUserByID(ctx context.Context, req *userspb.FindUserByIDRequest) (*userspb.User, error) {

	const op = "controller.grpc.FindUserByID"

	l := c.log.With(slog.String("op", op))

	metaParams := inputs.NewMetaParams()

	metaParams.SetLocale(grpchelper.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0])

	err := inputs.ValidateUUID(req.GetId(), metaParams.GetLocale())

	if err != nil {
		return nil, grpchelper.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	user, err := c.usecases.FindUserByID(ctx, req.GetId())

	if err != nil {
		return nil, handleError(l, metaParams.GetLocale(), keys.GetUserInternal, err)
	}

	return user.ToProto(), nil

}
