package worker

import (
	"fmt"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/goto/raccoon/collection"
	"github.com/goto/raccoon/logger"
	"github.com/goto/raccoon/metrics"
	"github.com/goto/raccoon/publisher"
)

// Pool spawn goroutine as much as Size that will listen to EventsChannel. On Close, wait for all data in EventsChannel to be processed.
type Pool struct {
	Size                int
	deliveryChannelSize int
	EventsChannel       <-chan collection.CollectRequest
	kafkaProducer       publisher.KafkaProducer
	wg                  sync.WaitGroup
}

// CreateWorkerPool create new Pool struct given size and EventsChannel worker.
func CreateWorkerPool(size int, eventsChannel <-chan collection.CollectRequest, deliveryChannelSize int, kafkaProducer publisher.KafkaProducer) *Pool {
	return &Pool{
		Size:                size,
		deliveryChannelSize: deliveryChannelSize,
		EventsChannel:       eventsChannel,
		kafkaProducer:       kafkaProducer,
		wg:                  sync.WaitGroup{},
	}
}

// StartWorkers initialize worker pool as much as Pool.Size
func (w *Pool) StartWorkers() {
	w.wg.Add(w.Size)
	for i := 0; i < w.Size; i++ {
		go func(workerName string) {
			logger.Info("Running worker: " + workerName)
			deliveryChan := make(chan kafka.Event, w.deliveryChannelSize)
			for request := range w.EventsChannel {
				metrics.Timing("batch_idle_in_channel_milliseconds", (time.Now().Sub(request.TimePushed)).Milliseconds(), "worker="+workerName)
				batchReadTime := time.Now()
				//@TODO - Should add integration tests to prove that the worker receives the same message that it produced, on the delivery channel it created

				err := w.kafkaProducer.ProduceBulk(request.GetEvents(), request.ConnectionIdentifier.Group, deliveryChan)

				produceTime := time.Since(batchReadTime)
				metrics.Timing("kafka_producebulk_tt_ms", produceTime.Milliseconds(), "")

				if request.AckFunc != nil {
					request.AckFunc(err)
				}

				totalErr := 0
				if err != nil {
					for _, err := range err.(publisher.BulkError).Errors {
						if err != nil {
							logger.Errorf("[worker] Fail to publish message to kafka %v", err)
							totalErr++
						}
					}
				}
				lenBatch := int64(len(request.GetEvents()))
				logger.Debug(fmt.Sprintf("Success sending messages, %v", lenBatch-int64(totalErr)))
				if lenBatch > 0 {
					eventTimingMs := time.Since(request.GetSentTime().AsTime()).Milliseconds() / lenBatch
					metrics.Timing("event_processing_duration_milliseconds", eventTimingMs, fmt.Sprintf("conn_group=%s", request.ConnectionIdentifier.Group))
					now := time.Now()
					metrics.Timing("worker_processing_duration_milliseconds", (now.Sub(batchReadTime).Milliseconds())/lenBatch, "worker="+workerName)
					metrics.Timing("server_processing_latency_milliseconds", (now.Sub(request.TimeConsumed)).Milliseconds()/lenBatch, fmt.Sprintf("conn_group=%s", request.ConnectionIdentifier.Group))
				}
			}
			w.wg.Done()
		}(fmt.Sprintf("worker-%d", i))
	}
}

// FlushWithTimeOut waits for the workers to complete the pending the messages
// to be flushed to the publisher within a timeout.
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
