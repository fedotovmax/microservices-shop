package controller

import (
	"context"
	"errors"
	"log/slog"

	"github.com/fedotovmax/grpcutils/violations"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/internal/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcController struct {
	userspb.UnimplementedUserServiceServer
	log      *slog.Logger
	usecases Usecases
}

func NewGRPCController(log *slog.Logger, u Usecases) *grpcController {
	return &grpcController{
		log:      log,
		usecases: u,
	}
}

func (c *grpcController) CreateUser(ctx context.Context, req *userspb.CreateUserRequest) (*userspb.CreateUserResponse, error) {

	locale := utils.GetFromMetadata(ctx, keys.MetadataLocaleKey, keys.FallbackLocale)

	input := domain.NewCreateUserInput(req.GetEmail(), req.GetPassword())

	err := input.Validate(locale[0])

	if err != nil {
		var ve violations.ValidationErrors
		if errors.As(err, &ve) {
			fieldviolations := ve.ToRPCViolations()

			badRequest := &errdetails.BadRequest{
				FieldViolations: fieldviolations,
			}

			st := status.New(codes.InvalidArgument, "validation failed")

			withDetails, err := st.WithDetails(badRequest)

			if err != nil {
				return nil, st.Err()
			}
			return nil, withDetails.Err()
		}
	}

	userId, err := c.usecases.CreateUser(ctx, input)

	if err != nil {
		st := status.New(codes.Internal, "internal error when create new user")
		return nil, st.Err()
	}

	return &userspb.CreateUserResponse{
		Id: userId,
	}, nil
}
