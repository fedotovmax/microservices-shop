package grpcadapter

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"google.golang.org/grpc"
)

var ErrForceStoppedServer = errors.New("the server was forcibly stopped due to a timeout")

type Config struct {
	Addr string
}

type Server struct {
	addr                  string
	usersSvc              userspb.UserServiceServer
	sessionActionSvc      userspb.SessionActionServiceServer
	verificationActionSvc userspb.VerificationServiceServer
	grpc                  *grpc.Server
}

func New(
	cfg Config,
	usersSvc userspb.UserServiceServer,
	sessionActionSvc userspb.SessionActionServiceServer,
	verificationActionSvc userspb.VerificationServiceServer,
	opt ...grpc.ServerOption,
) *Server {
	return &Server{
		addr:                  cfg.Addr,
		usersSvc:              usersSvc,
		sessionActionSvc:      sessionActionSvc,
		verificationActionSvc: verificationActionSvc,
		grpc:                  grpc.NewServer(opt...),
	}
}

func (s *Server) Start() error {

	const op = "adapter.grpc.Start"

	listener, err := net.Listen("tcp", s.addr)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	userspb.RegisterUserServiceServer(s.grpc, s.usersSvc)
	userspb.RegisterSessionActionServiceServer(s.grpc, s.sessionActionSvc)
	userspb.RegisterVerificationServiceServer(s.grpc, s.verificationActionSvc)

	if err := s.grpc.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	const op = "adapter.grpc.Stop2"

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
