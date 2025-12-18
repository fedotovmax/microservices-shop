package grpccontroller

import (
	"context"
	"log/slog"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/pkg/utils/grpchelper"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Usecases interface {
	CreateUser(ctx context.Context, d *inputs.CreateUserInput, emailIn *inputs.EmailVerifyNotificationInput) (string, error)
	UpdateUserProfile(ctx context.Context, id string, in *inputs.UpdateUserInput) error
	FindUserByID(ctx context.Context, id string) (*domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type grpcController struct {
	userspb.UnimplementedUserServiceServer
	log      *slog.Logger
	usecases Usecases
}

const (
	validationFailed = "validation failed"

	getUserInternal           = "internal error when get user"
	createUserInternal        = "internal error when create new user"
	updateUserProfileInternal = "internal error when update user profile"
)

func NewGRPCController(log *slog.Logger, u Usecases) *grpcController {
	return &grpcController{
		log:      log,
		usecases: u,
	}
}

// func (c *grpcController) FindUserByEmail(context.Context, *userspb.FindUserByEmailRequest) (*userspb.User, error)

// 	{

// 	}

func (c *grpcController) FindUserByID(ctx context.Context, req *userspb.FindUserByIDRequest) (*userspb.User, error) {

	const op = "controller.grpc.FindUserByID"

	l := c.log.With(slog.String("op", op))

	locale := grpchelper.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)

	err := inputs.ValidateUUID(req.GetId(), locale[0])

	if err != nil {
		return nil, grpchelper.ReturnGRPCBadRequest(l, validationFailed, err)
	}

	user, err := c.usecases.FindUserByID(ctx, req.GetId())

	if err != nil {
		return nil, handleError(l, locale[0], getUserInternal, err)
	}

	return user.ToProto(locale[0]), nil

}

func (c *grpcController) UpdateUserProfile(ctx context.Context, req *userspb.UpdateUserProfileRequest) (*emptypb.Empty, error) {

	const op = "controller.grpc.UpdateUserProfile"

	l := c.log.With(slog.String("op", op))

	locale := grpchelper.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)

	userID := grpchelper.GetFromMetadata(ctx, keys.MetadataUserIDKey, "")

	err := inputs.ValidateUUID(userID[0], locale[0])

	if err != nil {
		return nil, grpchelper.ReturnGRPCBadRequest(l, validationFailed, err)
	}

	updateUserProfileInput := inputs.NewUpdateUserInput()
	updateUserProfileInput.SetFromProto(req)

	err = updateUserProfileInput.Validate(locale[0])

	if err != nil {
		return nil, grpchelper.ReturnGRPCBadRequest(l, validationFailed, err)
	}

	err = c.usecases.UpdateUserProfile(ctx, userID[0], updateUserProfileInput)

	if err != nil {
		return nil, handleError(l, locale[0], updateUserProfileInternal, err)
	}

	return &emptypb.Empty{}, nil
}

func (c *grpcController) CreateUser(ctx context.Context, req *userspb.CreateUserRequest) (*userspb.CreateUserResponse, error) {

	const op = "controller.grpc.CreateUser"

	l := c.log.With(slog.String("op", op))

	locale := grpchelper.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)

	createUserInput := inputs.NewCreateUserInput()
	createUserInput.SetEmail(req.GetEmail())
	createUserInput.SetPassword(req.GetPassword())

	err := createUserInput.Validate(locale[0])

	if err != nil {
		return nil, grpchelper.ReturnGRPCBadRequest(l, validationFailed, err)
	}

	emailInput := inputs.NewEmailVerifyNotificationInput()
	emailInput.SetDescription("Текст описания письма")
	emailInput.SetTitle("Текст заголовка письма")
	emailInput.SetLocale(locale[0])

	userId, err := c.usecases.CreateUser(ctx, createUserInput, emailInput)

	if err != nil {
		return nil, handleError(l, locale[0], createUserInternal, err)
	}

	return &userspb.CreateUserResponse{
		Id: userId,
	}, nil
}
