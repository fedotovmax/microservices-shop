package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/grpcutils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *controller) SendNewEmailVerifyLink(ctx context.Context, req *userspb.SendNewEmailVerifyLinkRequest) (*emptypb.Empty, error) {

	const op = "controller.grpc.SendNewEmailVerifyLink"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	input := inputs.NewUUIDInput()
	input.SetUUID(req.GetUserId())

	err := input.Validate(locale, "UserID")

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	err = c.usecases.SendNewEmailVerifyLink(ctx, input.GetUUID(), locale)

	if err != nil {
		return nil, c.handleError(locale, keys.UpdateUserProfileInternal, err)
	}

	return &emptypb.Empty{}, nil

}
