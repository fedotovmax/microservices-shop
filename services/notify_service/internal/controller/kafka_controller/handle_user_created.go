package kafkacontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/microservices-shop-protos/events"
)

func (k *kafkaController) handleEmailVerifyLinkAdded(ctx context.Context, eventID string, payload []byte) error {
	const op = "controller.kafka_consumer.handleEmailVerifyLinkAdded"

	l := k.log.With(slog.String("op", op))

	err := k.usecases.IsNewEvent(ctx, eventID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var userEmailVerifyLinkAddedPayload events.UserEmailVerifyLinkAdded
	err = json.Unmarshal(payload, &userEmailVerifyLinkAddedPayload)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, ErrInvalidPayloadForEventType, err)
	}

	link := k.cfg.CustomerSiteURL + k.cfg.CustomerSiteURLEmailVerifyPath + fmt.Sprintf("/%s", userEmailVerifyLinkAddedPayload.EmailVerifyLink)

	l.Info("send email verification link to email", slog.String("link", link))

	return nil
}
