package users

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
)

type sendEmalVerifyLinkAddedEventParams struct {
	ID            string
	Email         string
	Link          string
	LinkExpiresAt time.Time
	Locale        string
}

func (u *usecases) SendEmalVerifyLinkAddedEvent(
	ctx context.Context,
	params *sendEmalVerifyLinkAddedEventParams,
) error {

	const op = "usecases.users.SendEmalVerifyLinkAddedEvent"

	linkAddedPayload := events.UserEmailVerifyLinkAdded{
		ID:                       params.ID,
		Email:                    params.Email,
		EmailVerifyLink:          params.Link,
		EmailVerifyLinkExpiresAt: params.LinkExpiresAt,
		Locale:                   params.Locale,
	}

	linkAddedPayloadBytes, err := json.Marshal(linkAddedPayload)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	linkAddedEventIn := inputs.NewCreateEventInput()
	linkAddedEventIn.SetAggregateID(params.ID)
	linkAddedEventIn.SetTopic(events.USER_EVENTS)
	linkAddedEventIn.SetType(events.USER_EMAIL_VERIFY_LINK_ADDED)
	linkAddedEventIn.SetPayload(linkAddedPayloadBytes)

	_, err = u.eventSender.CreateEvent(ctx, linkAddedEventIn)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
