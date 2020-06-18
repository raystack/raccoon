package worker

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"raccoon/logger"
	"raccoon/publisher"
	"sync"
)

// Pool spawn goroutine as much as Size that will listen to EventsChannel. On Close, wait for all data in EventsChannel to be processed.
type Pool struct {
	Size          int
	EventsChannel <-chan []byte
	kafkaProducer publisher.KafkaProducer
	wg            sync.WaitGroup
}

// CreateWorkerPool create new Pool struct given size and EventsChannel worker.
func CreateWorkerPool(size int, eventsChannel <-chan []byte, kafkaProducer publisher.KafkaProducer) *Pool {
	return &Pool{
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
				//@TODO - Should add integration tests to prove that the worker receives the same message that it produced, on the delivery channel it created
				if err := w.kafkaProducer.Produce(events, deliveryChan); err != nil {
					logger.Error(fmt.Sprintf("[worker] Fail to publish message to kafka %v", err))
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
