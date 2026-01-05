package grpcadapter

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcUsersClient struct {
	RPC          userspb.UserServiceClient
	conn         *grpc.ClientConn
	closeChannel chan struct{}
}

func NewUsersClient(addr string) (*grpcUsersClient, error) {

	const op = "client.grpc.user-client.NewUserClient"

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &grpcUsersClient{
		RPC:          userspb.NewUserServiceClient(conn),
		conn:         conn,
		closeChannel: make(chan struct{}),
	}, nil

}

func (c *grpcUsersClient) Stop(ctx context.Context) error {
	const op = "client.grpc.user-client.Stop"

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
