package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/pkg/utils/grpchelper"
)

func (c *grpcController) UserSessionAction(ctx context.Context, req *userspb.UserSessionActionRequest) (*userspb.UserSessionActionResponse, error) {

	const op = "controller.grpc.UserSessionAction"

	l := c.log.With(slog.String("op", op))
	metaParams := inputs.NewMetaParams()

	metaParams.SetLocale(grpchelper.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0])

	userSessionActionInput := inputs.NewSessionActionInput()
	userSessionActionInput.SetEmail(req.GetEmail())
	userSessionActionInput.SetPassword(req.GetPassword())

	err := userSessionActionInput.Validate(metaParams.GetLocale())

	if err != nil {
		return nil, grpchelper.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	sessionActionResponse, err := c.usecases.UserSessionAction(ctx, userSessionActionInput)

	if err != nil {
		return nil, handleError(l, metaParams.GetLocale(), keys.CreateUserInternal, err)
	}

	return sessionActionResponse.ToProto(), nil
}
