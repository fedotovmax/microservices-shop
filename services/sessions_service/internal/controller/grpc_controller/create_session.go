package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/grpcutils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/keys"
)

func (c *controller) CreateSession(ctx context.Context, req *sessionspb.CreateSessionRequest) (*sessionspb.CreateSessionResponse, error) {

	const op = "grpc_controller.CreateSession"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	bypassCode := grpcutils.GetFromMetadata(ctx, keys.MetadataBypassCodeKey, "")[0]

	input := inputs.NewPrepareSessionInput(
		req.Uid, req.UserAgent, req.Ip,
	)

	err := input.Validate(locale)

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	newSession, err := c.usecases.CreateSession(ctx, input, bypassCode)

	if err != nil {
		return nil, handleError(l, locale, keys.CreateSessionInternal, err)
	}

	return newSession.ToProto(), nil
}
