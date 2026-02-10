package sessions

import (
	"context"
	"fmt"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type clients struct {
	Sessions     sessionspb.SessionsServiceClient
	conn         *grpc.ClientConn
	closeChannel chan struct{}
}

func New(addr string) (*clients, error) {

	const op = "adapters.clients.grpc.sessions.New"

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &clients{
		Sessions:     sessionspb.NewSessionsServiceClient(conn),
		conn:         conn,
		closeChannel: make(chan struct{}),
	}, nil

}

func (c *clients) Stop(ctx context.Context) error {
	const op = "adapters.clients.grpc.sessions.Stop"

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
