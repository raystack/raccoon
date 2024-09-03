package worker

import (
	"sync"
	"testing"
	"time"

	"github.com/raystack/raccoon/collector"
	"github.com/raystack/raccoon/identification"
	pb "github.com/raystack/raccoon/proto"
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
		},
	}

	t.Run("StartWorkers", func(t *testing.T) {
		t.Run("Should publish messages on bufferChannel to kafka", func(t *testing.T) {
			kp := mockKafkaPublisher{}
			bc := make(chan collector.CollectRequest, 2)
			worker := Pool{
				Size:          1,
				EventsChannel: bc,
				producer:      &kp,
				wg:            sync.WaitGroup{},
			}
			worker.StartWorkers()

			kp.On("ProduceBulk", mock.Anything, mock.Anything, mock.Anything).Return(nil).Twice()
			kp.On("Name").Return("kafka")
			bc <- *request
			bc <- *request

			worker.FlushWithTimeOut(5 * time.Millisecond)

			kp.AssertExpectations(t)
		})
		t.Run("should call ack function", func(t *testing.T) {
			kp := mockKafkaPublisher{}
			kp.On("ProduceBulk", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			kp.On("Name").Return("kafka")
			defer kp.AssertExpectations(t)

			q := make(chan collector.CollectRequest, 1)
			worker := Pool{
				Size:          1,
				EventsChannel: q,
				producer:      &kp,
				wg:            sync.WaitGroup{},
			}
			worker.StartWorkers()

			ackMock := &mockAck{}
			ackMock.On("Ack", nil).Return().Once()
			defer ackMock.AssertExpectations(t)
			r := *request
			r.AckFunc = ackMock.Ack
			q <- r
			worker.FlushWithTimeOut(5 * time.Millisecond)
		})
	})

	t.Run("Flush", func(t *testing.T) {
		t.Run("Should block until all messages is processed", func(t *testing.T) {
			kp := mockKafkaPublisher{}
			bc := make(chan collector.CollectRequest, 2)

			worker := Pool{
				Size:          1,
				EventsChannel: bc,
				producer:      &kp,
				wg:            sync.WaitGroup{},
			}
			worker.StartWorkers()
			kp.On("ProduceBulk", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3).After(3 * time.Millisecond)
			kp.On("Name").Return("kafka")
			bc <- *request
			bc <- *request
			bc <- *request
			close(bc)
			timedOut := worker.FlushWithTimeOut(1 * time.Second)
			assert.False(t, timedOut)
			assert.Equal(t, 0, len(bc))
			kp.AssertExpectations(t)
		})
	})
}
