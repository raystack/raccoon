package worker

import (
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
}

// CreateWorkerPool create new Pool struct given size and EventsChannel worker.
func CreateWorkerPool(size int, eventsChannel <-chan []byte, kafkaProducer KafkaProducer) Pool {
	return Pool{
		Size:          size,
		EventsChannel: eventsChannel,
		kafkaProducer: kafkaProducer,
		wg:            sync.WaitGroup{},
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
					Value: events,
				}
				//@TODO - Should add integration tests to prove that the worker receives the same message that it produced, on the delivery channel it created
				onFailRetry(&message, deliveryChan, w.kafkaProducer)
			}
			w.wg.Done()
		}()
	}
}

// Flush wait for remaining data to be processed. Call this after closing EventsChannel channel
func (w *Pool) Flush() {
	w.wg.Wait()
}

func onFailRetry(message *kafka.Message, deliveryChan chan kafka.Event, producer KafkaProducer) {
	if err := producer.Produce(message, deliveryChan); err != nil {
		onFailRetry(message, deliveryChan, producer)
	}
}
