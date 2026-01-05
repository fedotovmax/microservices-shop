package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/grpcutils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/keys"
)

func (c *controller) VerifyAccessToken(ctx context.Context, req *sessionspb.VerifyAccessTokenRequest) (
	*sessionspb.VerifyAccessTokenResponse, error,
) {

	const op = "grpc_controller.VerifyAccessToken"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	input := inputs.NewVerifyAccessInput(
		req.AccessToken,
		req.Issuer,
	)

	err := input.Validate(locale)

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	session, err := c.usecases.VerifyAccessToken(ctx, input)

	if err != nil {
		return nil, handleError(l, locale, keys.VerifyAccessInternal, err)
	}

	return &sessionspb.VerifyAccessTokenResponse{
		Uid: session.User.Info.UID,
	}, nil

}
