package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop/notify_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/queries"
)

type IsNewEventUsecase struct {
	log   *slog.Logger
	query queries.Events
}

func NewIsNewEventUsecase(
	log *slog.Logger,
	query queries.Events,
) *IsNewEventUsecase {
	return &IsNewEventUsecase{
		log:   log,
		query: query,
	}
}

func (u *IsNewEventUsecase) Execute(ctx context.Context, eventID string) error {

	const op = "usecases.is_new_event"

	_, err := u.query.FindByID(ctx, eventID)

	if err != nil {
		if errors.Is(err, errs.ErrEventNotFound) {
			return nil
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return fmt.Errorf("%s: %w", op, errs.ErrEventAlreadyHandled)
}
