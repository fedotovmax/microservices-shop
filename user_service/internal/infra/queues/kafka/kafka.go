package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

func NewConsumerGroup(brokers []string, groupID string, topics []string) (sarama.ConsumerGroup, error) {

	const op = "kafka.NewConsumerGroup"

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_1_0_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.IsolationLevel = sarama.ReadUncommitted

	cg, err := sarama.NewConsumerGroup(brokers, groupID, cfg)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cg, nil
}

func NewAsyncProducer(brokers []string) (sarama.AsyncProducer, error) {
	const op = "kafka.NewAsyncProducer"
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Retry.Max = 5
	cfg.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return producer, nil
}
