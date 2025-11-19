package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

type Producer interface {
	GetInput() chan<- *sarama.ProducerMessage
	GetSuccesses() <-chan *sarama.ProducerMessage
	GetErrors() <-chan *sarama.ProducerError
	Close(context.Context) error
}

type producer struct {
	instance sarama.AsyncProducer
}

func (p *producer) GetInput() chan<- *sarama.ProducerMessage {
	return p.instance.Input()
}

func (p *producer) GetSuccesses() <-chan *sarama.ProducerMessage {
	return p.instance.Successes()
}

func (p *producer) GetErrors() <-chan *sarama.ProducerError {
	return p.instance.Errors()
}

// TODO: handle close with context
func (p *producer) Close(ctx context.Context) error {
	return p.instance.Close()
}

// TODO: sync once???
func NewAsyncProducer(brokers []string) (Producer, error) {
	const op = "queues.kafka.NewAsyncProducer"
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Retry.Max = 5
	cfg.Producer.RequiredAcks = sarama.WaitForAll

	p, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &producer{
		instance: p,
	}, nil
}
