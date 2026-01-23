package events

import (
	"context"
	"time"

	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/inputs"
	"github.com/fedotovmax/pgxtx"
)

type Storage interface {
	SetEventStatusDone(ctx context.Context, id string) error
	SetEventsReservedToByIDs(ctx context.Context, ids []string, dur time.Duration) error
	RemoveEventReserve(ctx context.Context, id string) error
	CreateEvent(ctx context.Context, d *inputs.CreateEvent) (string, error)
	FindNewAndNotReservedEvents(ctx context.Context, limit int) ([]*domain.Event, error)
}

type usecases struct {
	storage Storage
	txm     pgxtx.Manager
}

func New(storage Storage, txm pgxtx.Manager) *usecases {
	return &usecases{storage: storage, txm: txm}
}
