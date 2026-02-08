package eventspublisher

import (
	"context"

	"github.com/fedotovmax/kafka-lib/outbox"
	"github.com/fedotovmax/microservices-shop-protos/events"
)

type Publisher interface {
	SessionBlacklistAdded(ctx context.Context, payload events.SessionBlacklistAddedEventPayload) error
	SessionBypassAdded(ctx context.Context, payload events.SessionBypassAddedEventPayload) error
}

type publisher struct {
	eventSender outbox.Sender
}

func New(eventSender outbox.Sender) *publisher {
	return &publisher{
		eventSender: eventSender,
	}
}
