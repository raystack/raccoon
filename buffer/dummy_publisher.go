package buffer

import (
	"raccoon/logger"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type DummyPublisher struct{}

// TODO: Remove this once concrete producer available on master
func (d DummyPublisher) Produce(message *kafka.Message, deliveryChannel chan kafka.Event) error {
	deliveryChannel <- message
	time.Sleep(3 * time.Millisecond)
	<-deliveryChannel
	logger.Info("DELIVERED......")
	return nil
}
