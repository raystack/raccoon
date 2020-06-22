package publisher

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"raccoon/config"
	"raccoon/logger"
)

// KafkaProducer Produce data to kafka synchronously
type KafkaProducer interface {
	// ProduceBulk message to kafka. Block until all messages are sent. Return array of error. Order is not guaranteed.
	ProduceBulk(messages [][]byte, deliveryChannel chan kafka.Event) error
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

// ProduceBulk messages to kafka. Block until all messages are sent. Return array of error. Order of Errors is guaranteed.
// DeliveryChannel needs to be exclusive. DeliveryChannel is exposed for recyclability purpose.
func (pr *Kafka) ProduceBulk(data [][]byte, deliveryChannel chan kafka.Event) error {
	errors := make([]error, len(data))
	totalProcessed := 0
	for order, datum := range data {
		message := &kafka.Message{
			Value:          datum,
			TopicPartition: kafka.TopicPartition{Topic: &pr.Config.Topic, Partition: kafka.PartitionAny},
			Opaque:         order,
		}
		err := pr.kp.Produce(message, deliveryChannel)
		if err != nil {
			errors[order] = err
			continue
		}
		totalProcessed++
	}
	// Wait for deliveryChannel as many as processed
	for i := 0; i < totalProcessed; i++ {
		d := <-deliveryChannel
		m := d.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			order := m.Opaque.(int)
			errors[order] = m.TopicPartition.Error
		}
	}

	if allNil(errors) {
		return nil
	}
	return BulkError{Errors: errors}
}

// Close wait for outstanding messages to be delivered within given flush interval timeout.
func (pr *Kafka) Close() {
	remaining := pr.kp.Flush(pr.Config.GetFlushInterval())
	logger.Info(fmt.Sprintf("Total undelivered messages: %d", remaining))
	pr.kp.Close()
}

func allNil(errors []error) bool {
	for _, err := range errors {
		if err != nil {
			return false
		}
	}
	return true
}

type BulkError struct {
	Errors []error
}

func (b BulkError) Error() string {
	err := "error when sending messages: "
	for i, mErr := range b.Errors {
		if i != 0 {
			err += fmt.Sprintf(", %v", mErr)
			continue
		}
		err += mErr.Error()
	}
	return err
}
