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

func (c *session) CreateSession(ctx context.Context, req *sessionspb.CreateSessionRequest) (*sessionspb.CreateSessionResponse, error) {

	const op = "grpc_controller.session.CreateSession"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	input := inputs.NewPrepareSession(
		req.Uid, req.UserAgent, req.Ip, req.BypassCode, req.DeviceTrustToken,
	)

	err := input.Validate(locale)

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	newSession, err := c.createSession.Execute(ctx, input)

	if err != nil {
		return handleCreateSessionError(l, locale, keys.CreateSessionInternal, err)
	}

	var trustToken *sessionspb.CreatedTrustToken

	if newSession.TrustToken != nil {
		trustToken = &sessionspb.CreatedTrustToken{
			TrustTokenValue:   newSession.TrustToken.DeviceTrustTokenValue,
			TrustTokenExpTime: timestamppb.New(newSession.TrustToken.DeviceTrustTokenExpTime),
		}
	}

	return &sessionspb.CreateSessionResponse{
		Payload: &sessionspb.CreateSessionResponse_SessionCreated{
			SessionCreated: &sessionspb.SessionCreated{
				AccessToken:    newSession.AccessToken,
				RefreshToken:   newSession.RefreshToken,
				AccessExpTime:  timestamppb.New(newSession.AccessExpTime),
				RefreshExpTime: timestamppb.New(newSession.RefreshExpTime),
				TrustToken:     trustToken,
			},
		},
	}, nil

}
