package kafkacontroller

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/fedotovmax/kafka-lib/kafka"
	"github.com/fedotovmax/microservices-shop-protos/events"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/domain/errs"
	"github.com/fedotovmax/microservices-shop/sessions_service/internal/usecases"
	"github.com/fedotovmax/microservices-shop/sessions_service/pkg/logger"
)

var ErrKafkaMessagesChannelClosed = errors.New("messages channel was closed")
var ErrInvalidPayloadForEventType = errors.New("invalid payload for current event type")

type kafkaController struct {
	log        *slog.Logger
	createUser *usecases.CreateUserUsecase
}

func New(log *slog.Logger, createUser *usecases.CreateUserUsecase) *kafkaController {
	return &kafkaController{
		log:        log,
		createUser: createUser,
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

			commit := func() {
				s.MarkMessage(message, "")
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
				commit()
				continue
			}

			if eventType == "" {
				l.Error("empty event type")
				commit()
				continue
			}

			payload := message.Value

			l.Info("new event", slog.String("event_type", eventType), slog.String("event_id", eventID))

			switch eventType {

			case events.USER_CREATED:
				err := k.handleUserCreated(ctx, payload)
				if err != nil {
					k.handleErrors(err, commit, l)
					continue
				}
				commit()
			default:
				l.Error("invalid event type", slog.String("event_type", eventType))
				commit()
			}
		}
	}
}

func (k *kafkaController) handleErrors(err error, commit func(), l *slog.Logger) {

	switch {
	case errors.Is(err, ErrInvalidPayloadForEventType):
		l.Error("invalid payload", logger.Err(err))
		commit()
		return
	case errors.Is(err, errs.ErrInternalCreateUser):
		l.Error("Failed to create user, event will not be commited")
		return
	default:
		commit()
		return
	}
}
