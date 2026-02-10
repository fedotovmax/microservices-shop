package publisher

import (
	"context"

	"github.com/fedotovmax/kafka-lib/outbox"
)

type Publisher interface {
	UserEmalVerifyLinkAdded(ctx context.Context, params *UserEmalVerifyLinkAddedParams) error

	UserCreated(ctx context.Context, params *UserCreatedParams) error

	ProfileUpdated(ctx context.Context, params *ProfileUpdatedParams) error
}

type eventPublisher struct {
	creator outbox.Creator
}

func New(creator outbox.Creator) *eventPublisher {
	return &eventPublisher{
		creator: creator,
	}
}
