package usecase

import (
	"context"

	"github.com/fedotovmax/microservices-shop/users_service/internal/domain"
	"github.com/fedotovmax/outbox"
	"github.com/fedotovmax/pgxtx"
)

type Storage interface {
	Create(ctx context.Context, d *domain.CreateUserInput) (string, error)
}

type usecases struct {
	s   Storage
	es  outbox.EventSender
	txm pgxtx.Manager
}

func NewUsecases(s Storage, txm pgxtx.Manager, es outbox.EventSender) *usecases {
	return &usecases{
		s:   s,
		es:  es,
		txm: txm,
	}
}
