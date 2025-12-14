package controller

import (
	"context"
	"errors"
	"log/slog"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/pkg/utils/grpchelper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCUsecases interface {
	CreateUser(ctx context.Context, d *inputs.CreateUserInput) (string, error)
	UpdateUserProfile(ctx context.Context, id string, in *inputs.UpdateUserInput) error
	FindUserByID(ctx context.Context, id string) (*domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type grpcController struct {
	userspb.UnimplementedUserServiceServer
	log      *slog.Logger
	usecases GRPCUsecases
}

const (
	validationFailed = "validation failed"

	getUserInternal           = "internal error when get user"
	createUserInternal        = "internal error when create new user"
	updateUserProfileInternal = "internal error when update user profile"
)

func NewGRPCController(log *slog.Logger, u GRPCUsecases) *grpcController {
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
		if errors.Is(err, errs.ErrUserNotFound) {
			msg, err := i18n.Local.Get(locale[0], keys.UserNotFound)
			if err != nil {
				l.Error(err.Error())
			}
			st := status.New(codes.NotFound, msg)
			return nil, st.Err()
		}
		return nil, grpchelper.ReturnGRPCInternal(l, getUserInternal, err)
	}

	return user.ToProto(), nil

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

	if updateUserProfileInput.GetBirthDate() != nil {
		l.Info(*updateUserProfileInput.GetBirthDate())
	}

	err = c.usecases.UpdateUserProfile(ctx, userID[0], updateUserProfileInput)

	if err != nil {
		return nil, grpchelper.ReturnGRPCInternal(l, updateUserProfileInternal, err)
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

	userId, err := c.usecases.CreateUser(ctx, createUserInput)

	if err != nil {
		if errors.Is(err, errs.ErrUserAlreadyExists) {
			msg, err := i18n.Local.Get(locale[0], keys.UserAlreadyExists)
			if err != nil {
				l.Error(err.Error())
			}
			st := status.New(codes.AlreadyExists, msg)
			return nil, st.Err()
		}
		return nil, grpchelper.ReturnGRPCInternal(l, createUserInternal, err)
	}

	return &userspb.CreateUserResponse{
		Id: userId,
	}, nil
}
