package publisher

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"raccoon/config"
	"raccoon/logger"
)

// KafkaProducer Produce data to kafka synchronously
type KafkaProducer interface {
	Produce(message []byte, deliveryChannel chan kafka.Event) error
}

func NewKafka(config config.KafkaConfig) (*Kafka, error) {
	kp, err := newKafkaClient(config.ToKafkaConfigMap())
	if err != nil {
		return &Kafka{}, err
	}
	return &Kafka{
		kp:     kp,
		Config: config,
	}, nil
}

func NewKafkaFromClient(client Client, config config.KafkaConfig) *Kafka {
	return &Kafka{
		kp:     client,
		Config: config,
	}
}

type Kafka struct {
	kp     Client
	Config config.KafkaConfig
}

func (pr *Kafka) Produce(data []byte, deliveryChannel chan kafka.Event) error {
	message := &kafka.Message{
		Value:          data,
		TopicPartition: kafka.TopicPartition{Topic: &pr.Config.Topic, Partition: kafka.PartitionAny},
	}
	produceErr := pr.kp.Produce(message, deliveryChannel)

	if produceErr != nil {
		logger.Error("Producer failed to send message ", produceErr)
		return produceErr
	}

	e := <-deliveryChannel
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		logger.Error(fmt.Sprintf("Producer message delivery failed.%s", m.TopicPartition.Error))
		return m.TopicPartition.Error
	}
	logger.Debug(fmt.Sprintf("Delivered message to topic %s [%d] at offset %s",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset))
	return nil
}

// Close wait for outstanding messages to be delivered within given flush interval timeout.
func (pr *Kafka) Close() {
	remaining := pr.kp.Flush(pr.Config.GetFlushInterval())
	logger.Info(fmt.Sprintf("Total undelivered messages: %d", remaining))
	pr.kp.Close()
}
