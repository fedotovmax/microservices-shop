package grpcadapter

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/assortimentpb"
	"google.golang.org/grpc"
)

var ErrForceStoppedServer = errors.New("the server was forcibly stopped due to a timeout")

type Config struct {
	Addr string
}

type Server struct {
	addr        string
	brandSvc    assortimentpb.UnimplementedBrandServiceServer
	categorySvc assortimentpb.UnimplementedCategoryServiceServer
	grpc        *grpc.Server
}

func New(cfg Config,
	brandSvc assortimentpb.UnimplementedBrandServiceServer,
	categorySvc assortimentpb.UnimplementedCategoryServiceServer,
	opt ...grpc.ServerOption,
) *Server {
	return &Server{
		addr:        cfg.Addr,
		brandSvc:    brandSvc,
		categorySvc: categorySvc,
		grpc:        grpc.NewServer(opt...),
	}
}

func (s *Server) Start() error {

	const op = "adapters.grpc.Start"

	listener, err := net.Listen("tcp", s.addr)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	assortimentpb.RegisterBrandServiceServer(s.grpc, s.brandSvc)
	assortimentpb.RegisterCategoryServiceServer(s.grpc, s.categorySvc)

	if err := s.grpc.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	const op = "adapters.grpc.Stop"

	done := make(chan struct{})

	go func() {
		s.grpc.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		s.grpc.Stop()
		return fmt.Errorf("%s: %w", op, ErrForceStoppedServer)
	}
}
