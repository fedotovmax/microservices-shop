package usecases

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fedotovmax/kafka-lib/outbox"
	"github.com/fedotovmax/microservices-shop-protos/events"
)

type sendUserCreatedEventParams struct {
	ID     string
	Email  string
	Locale string
}

func (u *usecases) SendUserCreatedEvent(
	ctx context.Context,
	params *sendUserCreatedEventParams,
) error {

	const op = "usecases.users.SendUserCreatedEvent"

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

	_, err = u.eventSender.CreateEvent(ctx, userCreatedEventIn)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
