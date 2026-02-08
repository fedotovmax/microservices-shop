package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/grpcutils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/keys"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *session) RefreshSession(ctx context.Context, req *sessionspb.RefreshSessionRequest) (
	*sessionspb.CreateSessionResponse, error,
) {
	const op = "grpc_controller.session.RefreshSession"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	input := inputs.NewRefreshSession(
		req.RefreshToken, req.UserAgent, req.Ip,
	)

	err := input.Validate(locale)

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	response, err := c.refreshSession.Execute(ctx, input)

	if err != nil {
		return handleCreateSessionError(l, locale, keys.RefreshSessionInternal, err)
	}

	return &sessionspb.CreateSessionResponse{
		Payload: &sessionspb.CreateSessionResponse_SessionCreated{
			SessionCreated: &sessionspb.SessionCreated{
				AccessToken:    response.AccessToken,
				RefreshToken:   response.RefreshToken,
				AccessExpTime:  timestamppb.New(response.AccessExpTime),
				RefreshExpTime: timestamppb.New(response.RefreshExpTime),
			},
		},
	}, nil

}
