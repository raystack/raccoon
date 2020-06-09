package publisher

import (
	"clickstream-service/config"
	"clickstream-service/logger"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewProducer(kp KafkaProducer, config config.KafkaConfig) *Producer {
	return &Producer{
		kp:               kp,
		Config:           config,
	}
}

type Producer struct {
	kp               KafkaProducer
	InflightMessages chan *kafka.Message
	Config           config.KafkaConfig
}

func (pr *Producer) Produce(msg *kafka.Message) error {
	deliveryChan := make(chan kafka.Event)

	produceErr := pr.kp.Produce(msg, deliveryChan)

	if produceErr != nil {
		return produceErr
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		logger.Info(fmt.Sprintf("Kafka message delivery failed.%s",m.TopicPartition.Error))
	} else {
		 	logger.Info(fmt.Sprintf("Delivered message to topic %s [%d] at offset %s",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset))
		 	return m.TopicPartition.Error
	}

	close(deliveryChan)
	return nil
}

func (pr *Producer) Close() {
	pr.kp.Close()
}
