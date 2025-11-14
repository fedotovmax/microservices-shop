package usecase

import (
	"context"

	"github.com/fedotovmax/microservices-shop/user_service/internal/domain"
	"github.com/fedotovmax/microservices-shop/user_service/internal/ports"
)

type produceAdapter interface {
	Publish(ctx context.Context, ev domain.Event) error
	GetSuccesses(ctx context.Context) <-chan *domain.SuccessEvent
	GetErrors(ctx context.Context) <-chan *domain.FailedEvent
}

type eventProcessor struct {
	pa produceAdapter
	ea ports.EventAdapter
}

func NewEventProcessorUsecase(pa produceAdapter, ea ports.EventAdapter) *eventProcessor {
	return &eventProcessor{
		ea: ea,
		pa: pa,
	}
}

func (e *eventProcessor) DispatchMonitoring() {

}

func (e *eventProcessor) ProcessingNewEvents() {}
