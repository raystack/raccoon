package worker

import (
	"github.com/golang/protobuf/proto"
	"raccoon/logger"
	"source.golabs.io/mobile/clickstream-go-proto/gojek/clickstream/de"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"raccoon/publisher"
)

// Pool spawn goroutine as much as Size that will listen to EventsChannel. On Close, wait for all data in EventsChannel to be processed.
type Pool struct {
	Size                int
	deliveryChannelSize int
	EventsChannel       <-chan []*de.CSEventMessage
	kafkaProducer       publisher.KafkaProducer
	wg                  sync.WaitGroup
}

// CreateWorkerPool create new Pool struct given size and EventsChannel worker.
func CreateWorkerPool(size int, eventsChannel <-chan []*de.CSEventMessage, deliveryChannelSize int, kafkaProducer publisher.KafkaProducer) *Pool {
	return &Pool{
		Size:                size,
		deliveryChannelSize: deliveryChannelSize,
		EventsChannel:       eventsChannel,
		kafkaProducer:       kafkaProducer,
		wg:                  sync.WaitGroup{},
	}
}

// StartWorkers initialize worker pool as much as Pool.poolNumber
func (w *Pool) StartWorkers() {
	w.wg.Add(w.Size)
	for i := 0; i < w.Size; i++ {
		go func() {
			deliveryChan := make(chan kafka.Event, w.deliveryChannelSize)
			for events := range w.EventsChannel {
				//@TODO - Should add integration tests to prove that the worker receives the same message that it produced, on the delivery channel it created
				batch := make([][]byte, len(events))
				for _, event := range events {
					csByte, err := proto.Marshal(event)
					if err != nil {
						logger.Errorf("[worker] Fail to serialize message: %v", err)
						continue
					}
					batch = append(batch, csByte)
				}

				err := w.kafkaProducer.ProduceBulk(batch, deliveryChan)
				if err != nil {
					for _, err := range err.(publisher.BulkError).Errors {
						logger.Errorf("[worker] Fail to publish message to kafka %v", err)
					}
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
