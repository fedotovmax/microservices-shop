package ports

import (
	"context"
	"time"

	"github.com/fedotovmax/microservices-shop/user_service/internal/domain"
)

type EventAdapter interface {
	Create(ctx context.Context, d domain.CreateEvent) (string, error)
	SetReservedToByIDs(ctx context.Context, ids []string, dur time.Duration) error
	FindNewAndNotReserved(ctx context.Context, limit int) ([]*domain.Event, error)
	//RemoveReserveAndCnangeStatusByID(ctx context.Context, id string) error
	ChangeStatus(ctx context.Context, id string) error
	RemoveReserve(ctx context.Context, id string) error
}
