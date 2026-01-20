package grpccontroller

import (
	"context"
	"time"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/appsauthpb"
	"github.com/fedotovmax/microservices-shop/apps-auth/internal/domain"
	"github.com/fedotovmax/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *controller) CreateApp(ctx context.Context, req *appsauthpb.CreateAppRequest) (*appsauthpb.CreateAppResponse, error) {
	const op = "controller.CreateApp"

	name := req.GetName()

	err := validation.MinLength(name, 3)

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		return nil, st.Err()
	}

	appType := domain.ApplicationTypeFromProto(req.GetApplicationType())

	err = appType.IsValid()

	if err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		return nil, st.Err()
	}

	secret, err := c.usecases.CreateApp(ctx, &domain.App{CreatedAt: time.Now().UTC(), Name: name, Type: appType})

	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	return &appsauthpb.CreateAppResponse{Secret: secret}, nil

}
