package grpccontroller

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *controller) Test(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}
