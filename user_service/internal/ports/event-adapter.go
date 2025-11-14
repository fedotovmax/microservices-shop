package ports

import (
	"context"

	"github.com/fedotovmax/microservices-shop/user_service/internal/domain"
)

type EventAdapter interface {
	CreateEvent(ctx context.Context, d domain.CreateEvent) (string, error)
}
