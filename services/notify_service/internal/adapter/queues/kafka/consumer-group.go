package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

// TODO: implement
func NewConsumerGroup(cgcfg ConsumerGroupConfig) (sarama.ConsumerGroup, error) {
	const op = "queues.kafka.NewConsumerGroup"

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_1_0_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.IsolationLevel = sarama.ReadUncommitted

	cg, err := sarama.NewConsumerGroup(cgcfg.Brokers, cgcfg.GroupID, cfg)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cg, nil
}
