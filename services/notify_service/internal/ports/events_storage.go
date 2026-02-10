package ports

import "context"

type EventsStorage interface {
	FindByID(ctx context.Context, eventID string) (string, error)
	Save(ctx context.Context, eventID string) error
}
