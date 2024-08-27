package worker

import (
	"fmt"
	"sync"
	"time"

	"github.com/raystack/raccoon/collector"
	"github.com/raystack/raccoon/logger"
	"github.com/raystack/raccoon/metrics"
	pb "github.com/raystack/raccoon/proto"
	"github.com/raystack/raccoon/publisher"
)

// Producer produces data to sink
type Producer interface {
	// ProduceBulk message to a sink. Blocks until all messages are sent. Returns slice of error.
	ProduceBulk(events []*pb.Event, connGroup string) error

	// Name returns the name of the producer
	Name() string
}

// Pool spawn goroutine as much as Size that will listen to EventsChannel. On Close, wait for all data in EventsChannel to be processed.
type Pool struct {
	Size          int
	EventsChannel <-chan collector.CollectRequest
	producer      Producer
	wg            sync.WaitGroup
}

// CreateWorkerPool create new Pool struct given size and EventsChannel worker.
func CreateWorkerPool(size int, eventsChannel <-chan collector.CollectRequest, producer Producer) *Pool {
	return &Pool{
		Size:          size,
		EventsChannel: eventsChannel,
		producer:      producer,
		wg:            sync.WaitGroup{},
	}
}

func (w *Pool) worker(name string) {
	logger.Info("Running worker: " + name)
	for request := range w.EventsChannel {
		metrics.Histogram(
			"batch_idle_in_channel_milliseconds",
			time.Since(request.TimePushed).Milliseconds(),
			map[string]string{"worker": name})

		batchReadTime := time.Now()
		//@TODO - Should add integration tests to prove that the worker receives the same message that it produced, on the delivery channel it created

		err := w.producer.ProduceBulk(request.GetEvents(), request.ConnectionIdentifier.Group)

		produceTime := time.Since(batchReadTime)
		metrics.Histogram(
			fmt.Sprintf("%s_producebulk_tt_ms", w.producer.Name()),
			produceTime.Milliseconds(),
			map[string]string{},
		)

		if request.AckFunc != nil {
			request.AckFunc(err)
		}

		totalErr := 0
		if err != nil {
			switch et := err.(type) {
			case *publisher.BulkError:
				for _, e := range et.Errors {
					if e != nil {
						logger.Errorf("[worker] Fail to publish message: %v", e)
						totalErr++
					}
				}
			default:
				logger.Errorf("[worker] Failed to publish message: %v", et)
			}
		}

		lenBatch := int64(len(request.GetEvents()))
		logger.Debug(fmt.Sprintf("Success sending messages, %v", lenBatch-int64(totalErr)))
		if lenBatch > 0 {
			eventTimingMs := time.Since(request.GetSentTime().AsTime()).Milliseconds() / lenBatch
			metrics.Histogram(
				"event_processing_duration_milliseconds",
				eventTimingMs,
				map[string]string{"conn_group": request.ConnectionIdentifier.Group})
			now := time.Now()
			metrics.Histogram(
				"worker_processing_duration_milliseconds",
				(now.Sub(batchReadTime).Milliseconds())/lenBatch,
				map[string]string{"worker": name})
			metrics.Histogram(
				"server_processing_latency_milliseconds",
				(now.Sub(request.TimeConsumed)).Milliseconds()/lenBatch,
				map[string]string{"conn_group": request.ConnectionIdentifier.Group})
		}
	}
	w.wg.Done()
}

// StartWorkers initialize worker pool as much as Pool.Size
func (w *Pool) StartWorkers() {
	w.wg.Add(w.Size)
	for i := range w.Size {
		go w.worker(fmt.Sprintf("worker-%d", i))
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
