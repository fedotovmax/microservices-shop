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

func (c *controller) UpdateUserProfile(ctx context.Context, req *userspb.UpdateUserProfileRequest) (*emptypb.Empty, error) {

	const op = "controller.grpc.UpdateUserProfile"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	updateUserProfileInput := inputs.NewUpdateUserInput()

	updateUserProfileInput.SetFromProto(req)

	err := updateUserProfileInput.Validate(locale)

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	err = c.usecases.UpdateUserProfile(ctx, updateUserProfileInput, locale)

	if err != nil {
		return nil, c.handleError(l, locale, keys.UpdateUserProfileInternal, err)
	}

	return &emptypb.Empty{}, nil
}
