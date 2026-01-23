package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/grpcutils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/keys"
)

func (c *controller) RefreshSession(ctx context.Context, req *sessionspb.RefreshSessionRequest) (
	*sessionspb.CreateSessionResponse, error,
) {
	const op = "grpc_controller.RefreshSession"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	input := inputs.NewRefreshSessionInput(
		req.RefreshToken, req.UserAgent, req.Ip,
	)

	err := input.Validate(locale)

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	response, err := c.usecases.RefreshSession(ctx, input)

	if err != nil {
		return nil, handleError(l, locale, keys.RefreshSessionInternal, err)
	}

	return response.ToProto(), nil

}
