package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fedotovmax/kafka-lib/outbox"
	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
)

type ProfileUpdatedParams struct {
	UserID string
	Email  string
	Locale string
	Input  *inputs.UpdateUser
}

func (p *eventPublisher) ProfileUpdated(
	ctx context.Context,
	params *ProfileUpdatedParams,
) error {

	const op = "events_publisher.ProfileUpdated"

	userProfileUpdatedPayload := events.UserProfileUpdatedEventPayload{
		ID:            params.UserID,
		Email:         params.Email,
		NewLastName:   params.Input.GetLastName(),
		NewFirstName:  params.Input.GetFirstName(),
		NewMiddleName: params.Input.GetMiddleName(),
		NewAvatarURL:  params.Input.GetAvatarURL(),
		Locale:        params.Locale,
	}

	userProfileUpdatedPayloadBytes, err := json.Marshal(userProfileUpdatedPayload)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	userProfileUpdatedIn := outbox.NewCreateEventInput()
	userProfileUpdatedIn.SetAggregateID(params.UserID)
	userProfileUpdatedIn.SetTopic(events.USER_EVENTS)
	userProfileUpdatedIn.SetType(events.USER_PROFILE_UPDATED)
	userProfileUpdatedIn.SetPayload(userProfileUpdatedPayloadBytes)

	_, err = p.creator.CreateEvent(ctx, userProfileUpdatedIn)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
