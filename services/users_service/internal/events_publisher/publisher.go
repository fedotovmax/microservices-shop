package eventspublisher

import (
	"context"

	"github.com/fedotovmax/kafka-lib/outbox"
)

type Publisher interface {
	UserEmalVerifyLinkAdded(ctx context.Context, params *UserEmalVerifyLinkAddedParams) error

	UserCreated(ctx context.Context, params *UserCreatedParams) error

	ProfileUpdated(ctx context.Context, params *ProfileUpdatedParams) error
}

type publisher struct {
	eventSender outbox.Sender
}

func New(eventSender outbox.Sender) *publisher {
	return &publisher{
		eventSender: eventSender,
	}
}
