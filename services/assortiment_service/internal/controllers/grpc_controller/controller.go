package grpccontroller

import (
	"log/slog"

	"github.com/fedotovmax/microservices-shop-protos/gen/go/assortimentpb"
)

type Usecases interface{}

type controller struct {
	assortimentpb.UnimplementedAssortimentServiceServer
	log      *slog.Logger
	usecases Usecases
}

func New(log *slog.Logger, u Usecases) *controller {
	return &controller{
		log:      log,
		usecases: u,
	}
}
