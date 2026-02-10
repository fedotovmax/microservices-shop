package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/ports"
)

type SaveEventUsecase struct {
	log           *slog.Logger
	eventsStorage ports.EventsStorage
}

func NewSaveEventUsecase(
	log *slog.Logger,
	eventsStorage ports.EventsStorage,
) *SaveEventUsecase {
	return &SaveEventUsecase{
		log:           log,
		eventsStorage: eventsStorage,
	}
}

func (u *SaveEventUsecase) Execute(ctx context.Context, eventID string) error {

	const op = "usecases.save_event"

	err := u.eventsStorage.Save(ctx, eventID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
