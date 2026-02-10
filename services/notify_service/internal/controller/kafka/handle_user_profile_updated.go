package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/i18n"
	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/notify_service/internal/keys"
)

func (k *kafkaController) handleUserProfileUpdated(ctx context.Context, eventID string, payload []byte) error {

	const op = "controller.kafka_consumer.handleUserProfileUpdated"

	l := k.log.With(slog.String("op", op))

	err := k.isNewEvent.Execute(ctx, eventID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var userProfileUpdatedPayload events.UserProfileUpdatedEventPayload
	err = json.Unmarshal(payload, &userProfileUpdatedPayload)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, ErrInvalidPayloadForEventType, err)
	}

	sendCtx, cancelSendCtx := context.WithTimeout(ctx, time.Second*3)
	defer cancelSendCtx()

	text, err := i18n.Local.Get(userProfileUpdatedPayload.Locale, keys.ProfileUpdatedText)

	if err != nil {
		l.Warn(err.Error())
	}

	err = k.sendTgMessage.Execute(sendCtx, text, userProfileUpdatedPayload.ID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	saveEventCtx, cancelSaveEventCtx := context.WithTimeout(ctx, time.Second*2)
	defer cancelSaveEventCtx()

	err = k.saveEvent.Execute(saveEventCtx, eventID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
