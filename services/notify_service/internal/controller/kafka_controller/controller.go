package kafkacontroller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/fedotovmax/kafka-lib/kafka"
	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/notify_service/pkg/logger"
)

var ErrKafkaMessagesChannelClosed = errors.New("messages channel was closed")
var ErrInvalidPayloadForEventType = errors.New("invalid payload for current event type")

// TODO: real methods
type Usecases interface {
	SendTgMessage(ctx context.Context, text string, userId string) error
}

type kafkaController struct {
	log      *slog.Logger
	usecases Usecases
}

func NewKafkaController(log *slog.Logger, usecases Usecases) *kafkaController {
	return &kafkaController{
		log:      log,
		usecases: usecases,
	}
}

func (k *kafkaController) Setup(s sarama.ConsumerGroupSession) error {

	const op = "controller.kafka_consumer.Setup"

	l := k.log.With(slog.String("op", op))

	claims := s.Claims()

	l.Info("Setup: partitions assigned", slog.Any("claims", claims))

	// Examples for future maybe:
	// 1) Preload caches
	// k.cache.Load()

	// 2) Init workers for each partition
	// k.startWorkersForClaims(claims)

	// 3) Reset metrics
	// k.metrics.Reset()

	return nil
}

func (k *kafkaController) Cleanup(s sarama.ConsumerGroupSession) error {

	const op = "controller.kafka_consumer.Cleanup"

	l := k.log.With(slog.String("op", op))

	claims := s.Claims()

	l.Info("Cleanup: partitions revoked", slog.Any("claims", claims))

	// Examples for future maybe:
	// 1) Stop workers
	// k.stopWorkers()

	// 2) Flush buffers
	// k.flush()

	// 3) Close resources

	return nil

}

func (k *kafkaController) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {

	const op = "controller.kafka_consumer.ConsumeClaim"

	l := k.log.With(slog.String("op", op))

	ctx := s.Context()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%s: %w: %v", op, kafka.ErrConsumerHandlerClosedByCtx, ctx.Err())
		case message, ok := <-c.Messages():

			if !ok {
				return fmt.Errorf("%s: %w", op, ErrKafkaMessagesChannelClosed)
			}

			var eventID string
			var eventType string

			for _, header := range message.Headers {
				key := string(header.Key)
				switch key {
				case kafka.HeaderEventType:
					eventType = string(header.Value)
				case kafka.HeaderEventID:
					eventID = string(header.Value)
				}
			}

			if eventID == "" {
				l.Error("empty event ID")
				s.MarkMessage(message, "")
				continue
			}

			if eventType == "" {
				l.Error("empty event type")
				s.MarkMessage(message, "")
				continue
			}

			payload := message.Value

			switch eventType {

			case events.NOTIFICATIONS_EMAIL:

				var emailVerifyNotificationPayload events.EmailVerifyNotificationPayload

				err := json.Unmarshal(payload, &emailVerifyNotificationPayload)

				if err != nil {

					l.Error("invalid payload", logger.Err(err), slog.String("event_type", eventType))
					s.MarkMessage(message, "")
					continue
				}

				//TODO: real handle
				l.Info(
					"notify service: successfully consume message",
					slog.Any("payload", emailVerifyNotificationPayload),
					slog.Any("partition", message.Partition),
					slog.Int64("offset", message.Offset),
				)

				s.MarkMessage(message, "")

			case events.NOTIFICATIONS_TELEGRAM:

				err := k.handleTgNotification(ctx, payload)

				if err != nil {
					if errors.Is(err, ErrInvalidPayloadForEventType) {
						l.Error("invalid payload", logger.Err(err), slog.String("event_type", eventType))
						s.MarkMessage(message, "")
						continue
					}
					//todo: store events for retry maybe?
					l.Error("error when try to send message", logger.Err(err), slog.String("event_type", eventType))
					//todo: no mark, will retry via kafka offsets
					//s.MarkMessage(message, "")
					continue
				}

				s.MarkMessage(message, "")

			default:
				l.Error("invalid event type", slog.String("event_type", eventType))
				s.MarkMessage(message, "")
			}
		}
	}
}

func (k *kafkaController) handleTgNotification(ctx context.Context, payload []byte) error {

	const op = "controller.kafka_consumer.handleTgNotification"

	var tgNotificationPayload events.TelegramNotificationPayload
	err := json.Unmarshal(payload, &tgNotificationPayload)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, ErrInvalidPayloadForEventType, err)
	}

	sendCtx, cancelSetCtx := context.WithTimeout(ctx, time.Second*3)
	defer cancelSetCtx()

	err = k.usecases.SendTgMessage(sendCtx,
		tgNotificationPayload.Text, tgNotificationPayload.UserID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
