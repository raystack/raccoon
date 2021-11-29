package worker

import (
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"

	"raccoon/collection"
	"raccoon/identification"
	pb "raccoon/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorker(t *testing.T) {
	request := &collection.CollectRequest{
		ConnectionIdentifier: &identification.Identifier{
			ID:    "12345",
			Group: "viewer",
		},
		EventRequest: &pb.EventRequest{
			SentTime: &timestamp.Timestamp{Seconds: 1593574343},
		},
	}

	t.Run("StartWorkers", func(t *testing.T) {
		t.Run("Should publish messages on bufferChannel to kafka", func(t *testing.T) {
			kp := mockKafkaPublisher{}
			m := &mockMetric{}
			m.On("Timing", "processing.latency", mock.Anything, "")
			m.On("Count", "kafka_messages_delivered_total", 0, "success=true")
			m.On("Count", "kafka_messages_delivered_total", 0, "success=false")
			bc := make(chan *collection.CollectRequest, 2)
			worker := Pool{
				Size:                1,
				deliveryChannelSize: 0,
				EventsChannel:       bc,
				kafkaProducer:       &kp,
				wg:                  sync.WaitGroup{},
			}
			worker.StartWorkers()

			kp.On("ProduceBulk", mock.Anything, mock.Anything).Return(nil).Twice()
			bc <- request
			bc <- request
			time.Sleep(10 * time.Millisecond)

			kp.AssertExpectations(t)
		})
	})

	t.Run("Flush", func(t *testing.T) {
		t.Run("Should block until all messages is processed", func(t *testing.T) {
			kp := mockKafkaPublisher{}
			bc := make(chan *collection.CollectRequest, 2)
			m := &mockMetric{}
			m.On("Timing", "processing.latency", mock.Anything, "")
			m.On("Count", "kafka_messages_delivered_total", 0, "success=false")
			m.On("Count", "kafka_messages_delivered_total", 0, "success=true")

			worker := Pool{
				Size:                1,
				deliveryChannelSize: 100,
				EventsChannel:       bc,
				kafkaProducer:       &kp,
				wg:                  sync.WaitGroup{},
			}
			worker.StartWorkers()
			kp.On("ProduceBulk", mock.Anything, mock.Anything).Return(nil).Times(3).After(3 * time.Millisecond)
			bc <- request
			bc <- request
			bc <- request
			close(bc)
			timedOut := worker.FlushWithTimeOut(1 * time.Second)
			assert.False(t, timedOut)
			assert.Equal(t, 0, len(bc))
			kp.AssertExpectations(t)
		})
	})
}
