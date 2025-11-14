package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/IBM/sarama"
	"github.com/fedotovmax/microservices-shop/user_service/internal/domain"
	infraKafka "github.com/fedotovmax/microservices-shop/user_service/internal/infra/queues/kafka"
)

const producerClosed = "producer are closed"

const invalidMetadata = "invalid metadata"

type produceAdapter struct {
	producer infraKafka.Producer

	onceSuccess sync.Once
	onceErrors  sync.Once

	successes chan *domain.SuccessEvent
	errors    chan *domain.FailedEvent
}

func NewProduceAdapter(p infraKafka.Producer) *produceAdapter {
	return &produceAdapter{
		producer:  p,
		successes: make(chan *domain.SuccessEvent),
		errors:    make(chan *domain.FailedEvent),
	}
}

func (p *produceAdapter) Publish(ctx context.Context, ev domain.Event) error {
	op := "adapter.kafka.produce.Publish"

	metadata := &messageMetadata{
		ID:   ev.ID,
		Type: ev.Type,
	}

	msg := &sarama.ProducerMessage{
		Topic:    ev.Topic,
		Key:      sarama.StringEncoder(ev.AggregateID),
		Value:    sarama.ByteEncoder(ev.Payload),
		Metadata: metadata,
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", op, ctx.Err())
	case p.producer.GetInput() <- msg:
		return nil
	}
}

func (p *produceAdapter) GetSuccesses(ctx context.Context) <-chan *domain.SuccessEvent {

	const op = "adapter.kafka.produce.GetSuccesses"

	p.onceSuccess.Do(func() {
		go func() {
			defer close(p.successes)
			for {
				select {
				case <-ctx.Done():
					slog.Info("stopped by context:", slog.String("op", op))
					return
				case msg, ok := <-p.producer.GetSuccesses():
					if !ok {
						slog.Info(producerClosed, slog.String("op", op))
						return
					}
					m, ok := msg.Metadata.(*messageMetadata)
					if !ok {
						slog.Warn(invalidMetadata, slog.String("op", op), slog.Any("metadata", msg.Metadata))
						continue
					}
					select {
					case <-ctx.Done():
						slog.Info("stopped by context:", slog.String("op", op))
						return
					case p.successes <- &domain.SuccessEvent{ID: m.ID, Type: m.Type}:
					}
				}
			}
		}()
	})

	return p.successes
}

func (p *produceAdapter) GetErrors(ctx context.Context) <-chan *domain.FailedEvent {
	const op = "adapter.kafka.produce.GetErrors"

	p.onceErrors.Do(func() {
		go func() {
			defer close(p.errors)
			for {
				select {
				case <-ctx.Done():
					slog.Info("stopped by context:", slog.String("op", op))
					return
				case produceErr, ok := <-p.producer.GetErrors():
					if !ok {
						slog.Info(producerClosed, slog.String("op", op))
						return
					}

					m, ok := produceErr.Msg.Metadata.(*messageMetadata)
					if !ok {
						slog.Warn(invalidMetadata, slog.String("op", op), slog.Any("metadata", produceErr.Msg.Metadata))
						continue
					}
					select {
					case <-ctx.Done():
						slog.Info("stopped by context:", slog.String("op", op))
						return
					case p.errors <- &domain.FailedEvent{ID: m.ID, Type: m.Type, Error: produceErr.Err}:
					}
				}
			}
		}()
	})

	return p.errors
}
