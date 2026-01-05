package kafkacontroller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fedotovmax/microservices-shop-protos/events"
)

func (k *kafkaController) handleUserCreated(ctx context.Context, payload []byte) error {

	const op = "controller.kafka.handleUserCreated"

	var userCreatedEventPayload events.UserCreatedEventPayload
	err := json.Unmarshal(payload, &userCreatedEventPayload)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, ErrInvalidPayloadForEventType, err)
	}

	err = k.usecases.CreateUser(ctx, userCreatedEventPayload.ID, userCreatedEventPayload.Email)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
