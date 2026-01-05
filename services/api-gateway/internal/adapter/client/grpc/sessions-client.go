package grpcadapter

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcSessionsClient struct {
	RPC          sessionspb.SessionsServiceClient
	conn         *grpc.ClientConn
	closeChannel chan struct{}
}

func NewSessionsClient(addr string) (*grpcSessionsClient, error) {

	const op = "client.grpc.sessions-client.NewSessionsClient"

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &grpcSessionsClient{
		RPC:          sessionspb.NewSessionsServiceClient(conn),
		conn:         conn,
		closeChannel: make(chan struct{}),
	}, nil

}

func (c *grpcSessionsClient) Stop(ctx context.Context) error {
	const op = "client.grpc.sessions-client.Stop"

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
