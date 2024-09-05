package worker

import (
	"fmt"
	"testing"
	"time"

	"github.com/raystack/raccoon/clock"
	"github.com/raystack/raccoon/collector"
	"github.com/raystack/raccoon/identification"
	"github.com/raystack/raccoon/metrics"
	pb "github.com/raystack/raccoon/proto"
	"github.com/raystack/raccoon/publisher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestWorker(t *testing.T) {
	request := &collector.CollectRequest{
		ConnectionIdentifier: identification.Identifier{
			ID:    "12345",
			Group: "viewer",
		},
		SendEventRequest: &pb.SendEventRequest{
			SentTime: &timestamppb.Timestamp{},
			Events: []*pb.Event{
				{
					Type: "synthetic_event",
				},
			},
		},
	}

	t.Run("StartWorkers", func(t *testing.T) {
		t.Run("Should publish messages on bufferChannel to kafka", func(t *testing.T) {
			kp := mockKafkaPublisher{}
			kp.On("ProduceBulk", mock.Anything, mock.Anything, mock.Anything).Return(nil).Twice()
			kp.On("Name").Return("kafka")
			defer kp.AssertExpectations(t)

			bc := make(chan collector.CollectRequest, 2)
			worker := CreateWorkerPool(
				1, bc, &kp,
			)
			worker.StartWorkers()

			bc <- *request
			bc <- *request
			close(bc)

			assert.False(
				t,
				worker.FlushWithTimeOut(time.Second),
			)
		})
		t.Run("Should call ack function", func(t *testing.T) {
			kp := mockKafkaPublisher{}
			kp.On("ProduceBulk", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			kp.On("Name").Return("kafka")
			defer kp.AssertExpectations(t)

			q := make(chan collector.CollectRequest, 1)
			worker := CreateWorkerPool(
				1, q, &kp,
			)
			worker.StartWorkers()

			ackMock := &mockAck{}
			ackMock.On("Ack", nil).Return().Once()
			defer ackMock.AssertExpectations(t)
			req := *request
			req.AckFunc = ackMock.Ack
			q <- req
			close(q)
			assert.False(
				t,
				worker.FlushWithTimeOut(time.Second),
			)
		})
		t.Run("Should handle publisher error", func(t *testing.T) {

			e := &publisher.BulkError{
				Errors: []error{
					fmt.Errorf("simulated error"),
				},
			}
			kp := mockKafkaPublisher{}
			kp.On("ProduceBulk", mock.Anything, mock.Anything, mock.Anything).
				Return(e).
				Once()
			kp.On("ProduceBulk", mock.Anything, mock.Anything, mock.Anything).
				Return(fmt.Errorf("publisher failure")).
				Once()
			kp.On("Name").Return("kafka")
			defer kp.AssertExpectations(t)

			q := make(chan collector.CollectRequest, 2)
			worker := CreateWorkerPool(
				1, q, &kp,
			)
			worker.StartWorkers()
			q <- *request
			q <- *request
			close(q)
			assert.False(
				t,
				worker.FlushWithTimeOut(time.Second),
			)
		})
		t.Run("should publish metrics related to workers", func(t *testing.T) {
			kp := mockKafkaPublisher{}
			kp.On("ProduceBulk", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			kp.On("Name").Return("kafka")
			defer kp.AssertExpectations(t)

			eventsChannel := make(chan collector.CollectRequest, 1)

			now := time.Now()
			clk := &clock.Mock{}
			clk.On("Now").Return(now).Once()                           // batchReadTime
			clk.On("Now").Return(now.Add(2 * time.Millisecond)).Once() // produceTime
			clk.On("Now").Return(now.Add(3 * time.Millisecond)).Once() // eventTimingMs
			clk.On("Now").Return(now.Add(4 * time.Millisecond)).Once() // (worker_server)_processing_*
			defer clk.AssertExpectations(t)

			req := *request
			req.SentTime = timestamppb.New(now.Add(-3 * time.Millisecond)) // time when client sends the request
			req.TimeConsumed = now.Add(-2 * time.Millisecond)              // time when service handlers send request to collector
			req.TimePushed = now.Add(-time.Millisecond)                    // time when collector sends request to channel

			mockInstrument := &metrics.MockInstrument{}
			mockInstrument.On(
				"Histogram",
				"batch_idle_in_channel_milliseconds",
				int64(1),
				mock.Anything,
			).Return(nil).Once()
			mockInstrument.On(
				"Histogram",
				"kafka_producebulk_tt_ms",
				int64(2),
				mock.Anything,
			).Return(nil).Once()
			mockInstrument.On(
				"Histogram",
				"event_processing_duration_milliseconds",
				int64(6),
				mock.Anything,
			).Return(nil).Once()
			mockInstrument.On(
				"Histogram",
				"worker_processing_duration_milliseconds",
				int64(4),
				mock.Anything,
			).Return(nil).Once()
			mockInstrument.On(
				"Histogram",
				"server_processing_latency_milliseconds",
				int64(6),
				mock.Anything,
			).Return(nil).Once()

			defer mockInstrument.AssertExpectations(t)

			worker := &Pool{
				Size:          1,
				EventsChannel: eventsChannel,
				producer:      &kp,
				instrument:    mockInstrument,
				clock:         clk,
			}
			worker.StartWorkers()

			eventsChannel <- req
			close(eventsChannel)
			assert.False(
				t, worker.FlushWithTimeOut(time.Second),
			)
		})
	})

	t.Run("Flush", func(t *testing.T) {
		t.Run("Should block until all messages is processed", func(t *testing.T) {
			kp := mockKafkaPublisher{}
			bc := make(chan collector.CollectRequest, 2)

			worker := CreateWorkerPool(
				1, bc, &kp,
			)
			worker.StartWorkers()
			kp.On("ProduceBulk", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3).After(3 * time.Millisecond)
			kp.On("Name").Return("kafka")
			bc <- *request
			bc <- *request
			bc <- *request
			close(bc)
			timedOut := worker.FlushWithTimeOut(time.Second)
			assert.False(t, timedOut)
			assert.Equal(t, 0, len(bc))
			kp.AssertExpectations(t)
		})
	})
}
