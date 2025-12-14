package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/fedotovmax/kafka-lib/kafka"
	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/users_service/internal/keys"
	"github.com/fedotovmax/microservices-shop/users_service/pkg/logger"
)

var ErrKafkaMessagesChannelClosed = errors.New("messages channel was closed")

// TODO: real methods
type KafkaUsecases interface {
	Test(ctx context.Context, payload any)
}

type kafkaController struct {
	log      *slog.Logger
	usecases KafkaUsecases
}

func NewKafkaController(log *slog.Logger, usecases KafkaUsecases) *kafkaController {
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

	for {
		select {
		case <-s.Context().Done():
			return fmt.Errorf("%s: %w: %v", op, kafka.ErrConsumerHandlerClosedByCtx, s.Context().Err())
		case message, ok := <-c.Messages():

			if !ok {
				return fmt.Errorf("%s: %w", op, ErrKafkaMessagesChannelClosed)
			}

			var eventID string
			var eventType string

			for _, header := range message.Headers {
				key := string(header.Key)
				switch key {
				case keys.KafkaHeaderEventType:
					eventType = string(header.Value)
				case keys.KafkaHeaderEventID:
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

			case events.USER_CREATED:

				var createdUserPayload events.UserCreatedEventPayload

				err := json.Unmarshal(payload, &createdUserPayload)

				if err != nil {
					l.Error("invalid payload", logger.Err(err), slog.String("event_type", eventType))
					s.MarkMessage(message, "")
					continue
				}

				//TODO: real handle
				k.usecases.Test(s.Context(), createdUserPayload)
				l.Info(
					"======================successfully consume message",
					slog.Any("payload", createdUserPayload),
					slog.Any("partition", message.Partition),
					slog.Int64("offset", message.Offset),
				)

				s.MarkMessage(message, "")

			default:
				l.Error("invalid event type", slog.String("event_type", eventType))
				s.MarkMessage(message, "")
			}
		}
	}
}
