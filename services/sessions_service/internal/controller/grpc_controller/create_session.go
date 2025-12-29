package grpccontroller

import (
	"context"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/google/uuid"
)

func (c *controller) CreateSession(ctx context.Context, req *sessionspb.CreateSessionRequest) (*sessionspb.CreateSessionResponse, error) {
	const op = "grpc_controller.CreateSession"

	uid := uuid.New().String()
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36"

	ip := "45.59.125.130"

	issuer := "api-gateway-app"

	newSession, err := c.usecases.CreateSession(ctx, inputs.NewPrepareSessionInput(
		uid, ua, ip, issuer,
	))

	if err != nil {
		return nil, err
	}

	return newSession.ToProto(), nil
}
