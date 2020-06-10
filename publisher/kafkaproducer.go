package publisher

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"raccoon/config"
)

type KafkaProducer interface {
	Produce(*kafka.Message, chan kafka.Event) error
	Close()
	Flush(int) int
}

func NewKafkaProducer(cfg config.KafkaConfig) (KafkaProducer, error) {
	kp, err := kafka.NewProducer(cfg.ToKafkaConfigMap())
	if err != nil {
		return nil, err
	}
	return kp, nil
}
