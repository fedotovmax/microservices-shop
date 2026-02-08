package eventspublisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fedotovmax/kafka-lib/outbox"
	"github.com/fedotovmax/microservices-shop-protos/events"
)

func (p *publisher) SessionBypassAdded(ctx context.Context, payload events.SessionBypassAddedEventPayload) error {
	const op = "events_publisher.SessionBypassAdded"

	eventPayloadBytes, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	eventInput := outbox.NewCreateEventInput()
	eventInput.SetAggregateID(payload.UID)
	eventInput.SetTopic(events.SESSION_EVENTS)
	eventInput.SetType(events.SESSION_BYPASS_ADDED)
	eventInput.SetPayload(eventPayloadBytes)

	_, err = p.eventSender.CreateEvent(ctx, eventInput)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
