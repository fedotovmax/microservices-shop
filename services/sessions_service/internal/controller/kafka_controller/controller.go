package kafkacontroller

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/fedotovmax/kafka-lib/kafka"
)

var ErrKafkaMessagesChannelClosed = errors.New("messages channel was closed")

// TODO: real methods
type Usecases interface {
}

type kafkaController struct {
	log      *slog.Logger
	usecases Usecases
}

func New(log *slog.Logger, usecases Usecases) *kafkaController {
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
			//TODO:
			case "--------":

				_ = payload

			default:
				l.Error("invalid event type", slog.String("event_type", eventType))
				s.MarkMessage(message, "")
			}
		}
	}
}
