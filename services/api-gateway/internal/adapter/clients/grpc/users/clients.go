package users

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type clients struct {
	Users         userspb.UserServiceClient
	Verification  userspb.VerificationServiceClient
	SessionAction userspb.SessionActionServiceClient
	conn          *grpc.ClientConn
	closeChannel  chan struct{}
}

func New(addr string) (*clients, error) {

	const op = "adapters.clients.grpc.users.New"

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &clients{
		Users:         userspb.NewUserServiceClient(conn),
		Verification:  userspb.NewVerificationServiceClient(conn),
		SessionAction: userspb.NewSessionActionServiceClient(conn),
		conn:          conn,
		closeChannel:  make(chan struct{}),
	}, nil

}

func (c *clients) Stop(ctx context.Context) error {

	const op = "adapters.clients.grpc.users.Stop"

	done := make(chan error, 1)

	go func() {
		err := c.conn.Close()
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
