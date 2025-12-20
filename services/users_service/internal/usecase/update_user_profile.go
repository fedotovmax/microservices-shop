package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/users_service/internal/adapter"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
)

func (u *usecases) UpdateUserProfile(ctx context.Context, id string, in *inputs.UpdateUserInput) error {

	const op = "usecase.UpdateUserProfile"

	err := u.txm.Wrap(ctx, func(txCtx context.Context) error {

		user, err := u.FindUserByID(txCtx, id)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.s.UpdateUserProfile(txCtx, user.ID, in)

		if err != nil && !errors.Is(err, adapter.ErrNoFieldsToUpdate) {
			return fmt.Errorf("%s: %w", op, err)
		}

		//TODO: locale
		tgNotificationPayload := events.TelegramNotificationPayload{
			UserID: user.ID,
			Text:   "Ваш профиль был обновлен!",
		}
		tgNotificationPayloadBytes, err := json.Marshal(tgNotificationPayload)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		tgNotificationIn := inputs.NewCreateEventInput()
		tgNotificationIn.SetAggregateID(user.ID)
		tgNotificationIn.SetTopic(events.NOTIFICATIONS_EVENTS)
		tgNotificationIn.SetType(events.NOTIFICATIONS_TELEGRAM)
		tgNotificationIn.SetPayload(tgNotificationPayloadBytes)

		_, err = u.s.CreateEvent(txCtx, tgNotificationIn)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	return err
}
