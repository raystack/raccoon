package publisher

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Client interface {
	Produce(*kafka.Message, chan kafka.Event) error
	Close()
	Flush(int) int
	Events() chan kafka.Event
}

func newKafkaClient(cfg *kafka.ConfigMap) (Client, error) {
	kp, err := kafka.NewProducer(cfg)
	if err != nil {
		return nil, err
	}
	return kp, nil
}
