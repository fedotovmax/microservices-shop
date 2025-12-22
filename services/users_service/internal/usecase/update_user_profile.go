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

func (u *usecases) UpdateUserProfile(ctx context.Context, meta *inputs.MetaParams, in *inputs.UpdateUserInput) error {

	const op = "usecase.UpdateUserProfile"

	err := u.txm.Wrap(ctx, func(txCtx context.Context) error {

		user, err := u.FindUserByID(txCtx, meta.GetUserID())

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.s.UpdateUserProfile(txCtx, user.ID, in)

		if err != nil && !errors.Is(err, adapter.ErrNoFieldsToUpdate) {
			return fmt.Errorf("%s: %w", op, err)
		}

		//TODO: locale
		userProfileUpdatedPayload := events.UserUpdatedEventPayload{ID: user.ID, Locale: meta.GetLocale()}

		userProfileUpdatedPayloadBytes, err := json.Marshal(userProfileUpdatedPayload)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		userProfileUpdatedIn := inputs.NewCreateEventInput()
		userProfileUpdatedIn.SetAggregateID(user.ID)
		userProfileUpdatedIn.SetTopic(events.USER_EVENTS)
		userProfileUpdatedIn.SetType(events.USER_PROFILE_UPDATED)
		userProfileUpdatedIn.SetPayload(userProfileUpdatedPayloadBytes)

		_, err = u.s.CreateEvent(txCtx, userProfileUpdatedIn)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	return err
}
