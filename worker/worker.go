package worker

import (
	"fmt"
	"raccoon/logger"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaProducer Produce data to kafka synchronously
type KafkaProducer interface {
	Produce(message *kafka.Message, deliveryChannel chan kafka.Event) error
}

// Pool spawn goroutine as much as Size that will listen to EventsChannel. On Close, wait for all data in EventsChannel to be processed.
type Pool struct {
	Size          int
	EventsChannel <-chan []byte
	kafkaProducer KafkaProducer
	wg            sync.WaitGroup
	kafkaTopic    string
}

// CreateWorkerPool create new Pool struct given size and EventsChannel worker.
func CreateWorkerPool(size int, eventsChannel <-chan []byte, kafkaProducer KafkaProducer, kafkaTopic string) Pool {
	return Pool{
		Size:          size,
		EventsChannel: eventsChannel,
		kafkaProducer: kafkaProducer,
		wg:            sync.WaitGroup{},
		kafkaTopic:    kafkaTopic,
	}
}

// StartWorkers initialize worker pool as much as Pool.poolNumber
func (w *Pool) StartWorkers() {
	w.wg.Add(w.Size)
	for i := 0; i < w.Size; i++ {
		go func() {
			deliveryChan := make(chan kafka.Event, 1)
			for events := range w.EventsChannel {
				message := kafka.Message{
					Value:          events,
					TopicPartition: kafka.TopicPartition{Topic: &w.kafkaTopic, Partition: kafka.PartitionAny},
				}
				//@TODO - Should add integration tests to prove that the worker receives the same message that it produced, on the delivery channel it created
				if err := w.kafkaProducer.Produce(&message, deliveryChan); err != nil {
					logger.Info(fmt.Sprintf("[worker] Fail to publish message to kafka %v", err))
				}
			}
			w.wg.Done()
		}()
	}
}

// Flush wait for remaining data to be processed. Call this after closing EventsChannel channel
func (w *Pool) Flush() {
	w.wg.Wait()
}
