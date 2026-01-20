package grpccontroller

import (
	"context"
	"errors"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/appsauthpb"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/domain/errs"
	"github.com/fedotovmax/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *controller) GetToken(ctx context.Context, req *appsauthpb.GetTokenRequest) (*appsauthpb.GetTokenResponse, error) {
	const op = "controller.CreateApp"

	secret := req.GetSecret()

	err := validation.MinLength(secret, 3)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		return nil, st.Err()
	}

	token, err := c.usecases.GetToken(ctx, secret)

	if err != nil {
		if errors.Is(err, errs.ErrAppNotFound) {
			st := status.New(codes.NotFound, err.Error())
			return nil, st.Err()
		}
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	return &appsauthpb.GetTokenResponse{Token: token}, nil

}
