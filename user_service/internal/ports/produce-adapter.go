package ports

import (
	"context"

	"github.com/fedotovmax/microservices-shop/user_service/internal/domain"
)

type ProduceAdapter interface {
	Publish(ctx context.Context, ev *domain.Event) error
	GetSuccesses(ctx context.Context) <-chan *domain.SuccessEvent
	GetErrors(ctx context.Context) <-chan *domain.FailedEvent
}
