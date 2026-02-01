package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fedotovmax/kafka-lib/outbox"
	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/users_service/internal/adapters"
	"github.com/fedotovmax/microservices-shop/users_service/internal/domain/inputs"
)

func (u *usecases) UpdateUserProfile(ctx context.Context, in *inputs.UpdateUserInput, locale string) error {

	const op = "usecase.users.UpdateUserProfile"

	err := u.txm.Wrap(ctx, func(txCtx context.Context) error {

		user, err := u.FindUserByID(txCtx, in.GetUserID())

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.usersStorage.UpdateUserProfile(txCtx, user.ID, in)

		if err != nil && !errors.Is(err, adapters.ErrNoFieldsToUpdate) {
			return fmt.Errorf("%s: %w", op, err)
		}

		user, err = u.FindUserByID(txCtx, user.ID)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		userProfileUpdatedPayload := events.UserProfileUpdatedEventPayload{
			ID:            user.ID,
			Email:         user.Email,
			NewLastName:   user.Profile.LastName,
			NewFirstName:  user.Profile.FirstName,
			NewMiddleName: user.Profile.MiddleName,
			NewAvatarURL:  user.Profile.AvatarURL,
			Locale:        locale,
		}

		userProfileUpdatedPayloadBytes, err := json.Marshal(userProfileUpdatedPayload)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		userProfileUpdatedIn := outbox.NewCreateEventInput()
		userProfileUpdatedIn.SetAggregateID(user.ID)
		userProfileUpdatedIn.SetTopic(events.USER_EVENTS)
		userProfileUpdatedIn.SetType(events.USER_PROFILE_UPDATED)
		userProfileUpdatedIn.SetPayload(userProfileUpdatedPayloadBytes)

		_, err = u.eventSender.CreateEvent(txCtx, userProfileUpdatedIn)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	return err
}
