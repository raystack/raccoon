package worker

import (
	"fmt"
	"testing"
	"time"

	"github.com/raystack/raccoon/collector"
	"github.com/raystack/raccoon/identification"
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
			r := *request
			r.AckFunc = ackMock.Ack
			q <- r
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

			ackMock := &mockAck{}
			ackMock.On("Ack", mock.Anything).Return().Twice()
			defer ackMock.AssertExpectations(t)
			r := *request
			r.AckFunc = ackMock.Ack
			q <- r
			q <- r
			close(q)
			assert.False(
				t,
				worker.FlushWithTimeOut(time.Second),
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
