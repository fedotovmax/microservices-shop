package kafkacontroller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/fedotovmax/kafka-lib/kafka"
	"github.com/fedotovmax/microservices-shop-protos/events"
)

var ErrKafkaMessagesChannelClosed = errors.New("messages channel was closed")
var ErrInvalidPayloadForEventType = errors.New("invalid payload for current event type")

// TODO: real methods
type Usecases interface {
	SendTgMessage(ctx context.Context, text string, userId string) error
	SaveEvent(ctx context.Context, eventID string) error
	IsNewEvent(ctx context.Context, eventID string) error
}

type Config struct {
	CustomerSiteURL                string
	CustomerSiteURLEmailVerifyPath string
}

type kafkaController struct {
	log      *slog.Logger
	usecases Usecases
	cfg      *Config
}

func NewKafkaController(log *slog.Logger, usecases Usecases, cfg *Config) *kafkaController {
	return &kafkaController{
		log:      log,
		usecases: usecases,
		cfg:      cfg,
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

			case events.USER_PROFILE_UPDATED:
				err := k.handleUserProfileUpdated(ctx, eventID, payload)
				if err != nil {
					k.handleErrors(err, commit, l)
					continue
				}
				commit()
			case events.SESSION_BYPASS_ADDED:
				//TODO: handle
				var p events.SessionBypassAddedEventPayload
				err := json.Unmarshal(payload, &p)

				if err != nil {
					k.handleErrors(ErrInvalidPayloadForEventType, commit, l)
					continue
				}
				l.Info("Send Email", slog.String("email", p.Email), slog.String("code", p.Code))
				commit()
			case events.USER_EMAIL_VERIFY_LINK_ADDED:
				err := k.handleEmailVerifyLinkAdded(ctx, eventID, payload)
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
