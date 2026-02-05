package usecases

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fedotovmax/kafka-lib/outbox"
	"github.com/fedotovmax/microservices-shop-protos/events"
)

type createUserCreatedEventParams struct {
	ID     string
	Email  string
	Locale string
}

func (u *usecases) createUserCreatedEvent(
	ctx context.Context,
	params *createUserCreatedEventParams,
) error {

	const op = "usecases.users.createUserCreatedEvent"

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
