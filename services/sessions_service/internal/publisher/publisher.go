package publisher

import (
	"context"

	"github.com/fedotovmax/kafka-lib/outbox"
	"github.com/fedotovmax/microservices-shop-protos/events"
)

type Publisher interface {
	SessionBlacklistAdded(ctx context.Context, payload events.SessionBlacklistAddedEventPayload) error
	SessionBypassAdded(ctx context.Context, payload events.SessionBypassAddedEventPayload) error
}

type eventPublisher struct {
	creator outbox.Creator
}

func New(creator outbox.Creator) *eventPublisher {
	return &eventPublisher{
		creator: creator,
	}
}
