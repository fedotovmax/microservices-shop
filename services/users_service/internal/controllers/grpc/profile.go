package grpc

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/grpcutils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/internal/queries"
	"github.com/fedotovmax/microservices-shop/users_service/internal/usecases"
	"google.golang.org/protobuf/types/known/emptypb"
)

type profile struct {
	userspb.UnimplementedUserServiceServer
	log           *slog.Logger
	updateProfile *usecases.UpdateProfileUsecase
	createUser    *usecases.CreateUserUsecase
	query         queries.Users
}

func NewProfile(
	log *slog.Logger,
	updateProfile *usecases.UpdateProfileUsecase,
	createUser *usecases.CreateUserUsecase,
	query queries.Users,
) *profile {
	return &profile{
		log:           log,
		updateProfile: updateProfile,
		createUser:    createUser,
		query:         query,
	}
}

func (c *profile) UpdateUserProfile(ctx context.Context, req *userspb.UpdateUserProfileRequest) (*emptypb.Empty, error) {

	const op = "controller.grpc.profile.UpdateUserProfile"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	updateUserProfileInput := inputs.NewUpdateUser()

	updateUserProfileInput.SetFromProto(req)

	err := updateUserProfileInput.Validate(locale)

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	err = c.updateProfile.Execute(ctx, updateUserProfileInput, locale)

	if err != nil {
		return nil, handleError(l, locale, keys.UpdateUserProfileInternal, err)
	}

	return &emptypb.Empty{}, nil
}

func (c *profile) FindUserByID(ctx context.Context, req *userspb.FindUserByIDRequest) (*userspb.User, error) {

	const op = "controller.grpc.profile.FindUserByID"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	findUserByIdInput := inputs.NewUUIDInput()
	findUserByIdInput.SetUUID(req.GetId())

	err := findUserByIdInput.Validate(locale, "UserID")

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	user, err := c.query.FindByID(ctx, req.GetId())

	if err != nil {
		return nil, handleError(l, locale, keys.GetUserInternal, err)
	}

	return user.ToProto(), nil
}

func (c *profile) CreateUser(ctx context.Context, req *userspb.CreateUserRequest) (*userspb.CreateUserResponse, error) {

	const op = "controller.grpc.profile.CreateUser"

	l := c.log.With(slog.String("op", op))

	locale := grpcutils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)[0]

	createUserInput := inputs.NewCreateUser()
	createUserInput.SetEmail(req.GetEmail())
	createUserInput.SetPassword(req.GetPassword())

	err := createUserInput.Validate(locale)

	if err != nil {
		return nil, grpcutils.ReturnGRPCBadRequest(l, keys.ValidationFailed, err)
	}

	userId, err := c.createUser.Execute(ctx, createUserInput, locale)

	if err != nil {
		return nil, handleError(l, locale, keys.CreateUserInternal, err)
	}

	return &userspb.CreateUserResponse{
		Id: userId,
	}, nil
}
