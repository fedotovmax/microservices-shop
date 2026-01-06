package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/grpcutils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
)

func (c *controller) UserSessionAction(ctx context.Context, req *userspb.UserSessionActionRequest) (*userspb.UserSessionActionResponse, error) {

	const op = "controller.grpc.UserSessionAction"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	userSessionActionInput := inputs.NewSessionActionInput()
	userSessionActionInput.SetEmail(req.GetEmail())
	userSessionActionInput.SetPassword(req.GetPassword())

	err := userSessionActionInput.Validate(locale)

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	sessionActionResponse, err := c.usecases.UserSessionAction(ctx, userSessionActionInput)

	if err != nil {
		return nil, c.handleError(l, locale, keys.CreateUserInternal, err)
	}

	return sessionActionResponse.ToProto(), nil
}
