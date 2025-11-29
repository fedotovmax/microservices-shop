package httpserver

import (
	"context"
	"fmt"
	"net/http"
)

type HTTPServerConfig struct {
	Port uint16
}

type Server struct {
	http *http.Server
}

func NewHTTPServer(c HTTPServerConfig, h http.Handler) *Server {

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.Port),
		Handler: h,
	}

	return &Server{
		http: s,
	}
}

func (s *Server) Start() error {

	const op = "server.http.Start"

	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (c *Server) Stop(ctx context.Context) error {

	const op = "server.http.Stop"

	err := c.http.Shutdown(ctx)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
