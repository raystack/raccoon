package worker

import (
	"raccoon/logger"
	"raccoon/metrics"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"source.golabs.io/mobile/clickstream-go-proto/gojek/clickstream/de"

	"raccoon/publisher"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Pool spawn goroutine as much as Size that will listen to EventsChannel. On Close, wait for all data in EventsChannel to be processed.
type Pool struct {
	Size                int
	deliveryChannelSize int
	EventsChannel       <-chan de.EventRequest
	kafkaProducer       publisher.KafkaProducer
	wg                  sync.WaitGroup
	instrumentation     metric
}

// CreateWorkerPool create new Pool struct given size and EventsChannel worker.
func CreateWorkerPool(size int, eventsChannel <-chan de.EventRequest, deliveryChannelSize int, kafkaProducer publisher.KafkaProducer) *Pool {
	return &Pool{
		Size:                size,
		deliveryChannelSize: deliveryChannelSize,
		EventsChannel:       eventsChannel,
		kafkaProducer:       kafkaProducer,
		wg:                  sync.WaitGroup{},
		instrumentation:     metrics.Instance(),
	}
}

// StartWorkers initialize worker pool as much as Pool.Size
func (w *Pool) StartWorkers() {
	w.wg.Add(w.Size)
	for i := 0; i < w.Size; i++ {
		go func() {
			deliveryChan := make(chan kafka.Event, w.deliveryChannelSize)
			for request := range w.EventsChannel {
				//@TODO - Should add integration tests to prove that the worker receives the same message that it produced, on the delivery channel it created
				batch := make([][]byte, 0, len(request.GetData()))
				for _, event := range request.GetData() {
					csByte, err := proto.Marshal(event)
					if err != nil {
						logger.Errorf("[worker] Fail to serialize message: %v", err)
						w.instrumentation.Count("kafka.event.serialization.error", 1, "")
						break
					}
					batch = append(batch, csByte)
				}
				err := w.kafkaProducer.ProduceBulk(batch, deliveryChan)
				totalErr := 0
				if err != nil {
					for _, err := range err.(publisher.BulkError).Errors {
						if err != nil {
							logger.Errorf("[worker] Fail to publish message to kafka %v", err)
							totalErr++
						}
					}
				}
				logger.Infof("Success sending messages, %v", len(batch))
				if len(batch) > 0 {
					eventTimingMs := time.Since(time.Unix(request.SentTime.Seconds, 0)).Milliseconds() / int64(len(batch))
					w.instrumentation.Timing("processing.latency", eventTimingMs, "")
				}
				w.instrumentation.Count("kafka.messages.delivered", totalErr, "success=false")
				w.instrumentation.Count("kafka.messages.delivered", len(batch)-totalErr, "success=true")
			}
			w.wg.Done()
		}()
	}
}

// FlushWithTimeOut waits for the workers to complete the pending the messages
//to be flushed to the publisher within a timeout.
// Returns true if waiting timed out, meaning not all the events could be processed before this timeout.
func (w *Pool) FlushWithTimeOut(timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		w.wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

type metric interface {
	Count(string, int, string)
	Timing(string, int64, string)
}
