package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

// TODO: sync once??
func NewConsumerGroup(brokers []string, groupID string, topics []string) (sarama.ConsumerGroup, error) {
	const op = "queues.kafka.NewConsumerGroup"

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
