package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/pkg/utils/grpchelper"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *grpcController) UpdateUserProfile(ctx context.Context, req *userspb.UpdateUserProfileRequest) (*emptypb.Empty, error) {

	const op = "controller.grpc.UpdateUserProfile"

	l := c.log.With(slog.String("op", op))

	metaParams := inputs.NewMetaParams()

	metaParams.SetLocale(grpchelper.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0])

	userID := grpchelper.GetFromMetadata(ctx, keys.MetadataUserIDKey, "")[0]

	err := inputs.ValidateUUID(userID, metaParams.GetLocale())

	if err != nil {
		return nil, grpchelper.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	metaParams.SetUserID(userID)

	updateUserProfileInput := inputs.NewUpdateUserInput()
	updateUserProfileInput.SetFromProto(req)

	err = updateUserProfileInput.Validate(metaParams.GetLocale())

	if err != nil {
		return nil, grpchelper.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	err = c.usecases.UpdateUserProfile(ctx, metaParams, updateUserProfileInput)

	if err != nil {
		return nil, handleError(l, metaParams.GetLocale(), keys.UpdateUserProfileInternal, err)
	}

	return &emptypb.Empty{}, nil
}
