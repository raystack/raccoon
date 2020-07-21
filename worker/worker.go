package worker

import (
	"fmt"
	"raccoon/logger"
	"raccoon/metrics"
	ws "raccoon/websocket"
	"sync"
	"time"

	"raccoon/publisher"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Pool spawn goroutine as much as Size that will listen to EventsChannel. On Close, wait for all data in EventsChannel to be processed.
type Pool struct {
	Size                int
	deliveryChannelSize int
	EventsChannel       <-chan ws.EventsBatch
	kafkaProducer       publisher.KafkaProducer
	wg                  sync.WaitGroup
}

// CreateWorkerPool create new Pool struct given size and EventsChannel worker.
func CreateWorkerPool(size int, eventsChannel <-chan ws.EventsBatch, deliveryChannelSize int, kafkaProducer publisher.KafkaProducer) *Pool {
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
				metrics.Timing("batch.idletime.inchannel", (time.Now().Sub(request.TimePushed)).Milliseconds(), "worker="+workerName)
				batchReadTime := time.Now()
				//@TODO - Should add integration tests to prove that the worker receives the same message that it produced, on the delivery channel it created
				batch := make([][]byte, 0, len(request.EventReq.GetEvents()))
				for _, event := range request.EventReq.GetEvents() {
					csByte := event.GetEventBytes()
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
					eventTimingMs := time.Since(time.Unix(request.EventReq.SentTime.Seconds, 0)).Milliseconds() / int64(len(batch))
					logger.Info(fmt.Sprintf("Currenttim: %d, eventTimingMs: %d", request.EventReq.SentTime.Seconds, eventTimingMs))
					metrics.Timing("processing.latency", eventTimingMs, "")
					metrics.Timing("worker.processing.latency", (time.Now().Sub(batchReadTime).Milliseconds())/int64(len(batch)), "worker="+workerName)
				}
				metrics.Count("kafka.messages.delivered", totalErr, "success=false")
				metrics.Count("kafka.messages.delivered", len(batch)-totalErr, "success=true")
			}
			w.wg.Done()
		}(fmt.Sprintf("worker-%d", i))
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
