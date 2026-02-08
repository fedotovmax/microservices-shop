package eventspublisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fedotovmax/kafka-lib/outbox"
	"github.com/fedotovmax/microservices-shop-protos/events"
)

type UserCreatedParams struct {
	ID     string
	Email  string
	Locale string
}

func (p *publisher) UserCreated(
	ctx context.Context,
	params *UserCreatedParams,
) error {

	const op = "events_publisher.UserCreated"

	userCreatedPayload := events.UserCreatedEventPayload{
		ID:     params.ID,
		Email:  params.Email,
		Locale: params.Locale,
	}

	userCreatedPayloadBytes, err := json.Marshal(userCreatedPayload)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	userCreatedEventIn := outbox.NewCreateEventInput()
	userCreatedEventIn.SetAggregateID(params.ID)
	userCreatedEventIn.SetTopic(events.USER_EVENTS)
	userCreatedEventIn.SetType(events.USER_CREATED)
	userCreatedEventIn.SetPayload(userCreatedPayloadBytes)

	_, err = p.eventSender.CreateEvent(ctx, userCreatedEventIn)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
