package publisher

import (
	"clickstream-service/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer interface {
	Produce(*kafka.Message, chan kafka.Event) error
	Close()
}

func NewKafkaProducer(cfg config.KafkaConfig) (KafkaProducer, error) {
	kp, err := kafka.NewProducer(cfg.ToKafkaConfigMap())
	if err != nil {
		return nil, err
	}
	return kp, nil
}
