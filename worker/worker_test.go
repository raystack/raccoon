package worker

import (
	"errors"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"raccoon/publisher"
	ws "raccoon/websocket"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"

	"source.golabs.io/mobile/clickstream-go-proto/gojek/clickstream/de"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorker(t *testing.T) {
	request := ws.EventsBatch{
		EventReq: de.EventRequest{
			SentTime: &timestamp.Timestamp{Seconds: 1593574343},
		},
	}

	t.Run("StartWorkers", func(t *testing.T) {
		t.Run("Should publish messages on bufferChannel to kafka", func(t *testing.T) {
			kp := mockKafkaPublisher{}
			m := &mockMetric{}
			m.On("Timing", "processing.latency", mock.Anything, "")
			m.On("Count", "kafka.messages.delivered", 0, "success=true")
			m.On("Count", "kafka.messages.delivered", 0, "success=false")
			bc := make(chan ws.EventsBatch, 2)
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

		t.Run("Should skip event with error topic", func(t *testing.T) {
			kp := mockKafkaPublisher{}
			bc := make(chan ws.EventsBatch, 2)
			tc := mockTopicCreator{}
			tc.On("CreateTopics", mock.Anything, mock.Anything, mock.Anything).Return([]kafka.TopicResult{{}}, errors.New("error"))
			worker := Pool{
				Size:                1,
				deliveryChannelSize: 0,
				EventsChannel:       bc,
				kafkaProducer:       &kp,
				wg:                  sync.WaitGroup{},
				router: &Router{
					topicsCreator:     &tc,
					format:            "%s",
					numPartitions:     0,
					replicationFactor: 0,
					m:                 &sync.Mutex{},
					topics:            make(map[string]string),
				},
			}
			worker.StartWorkers()

			kp.On("ProduceBulk", make([]publisher.Event, len(request.EventReq.GetEvents())), mock.Anything).Return(nil).Once()
			bc <- request
			time.Sleep(10 * time.Millisecond)

			kp.AssertExpectations(t)
		})
	})

	t.Run("Flush", func(t *testing.T) {
		t.Run("Should block until all messages is processed", func(t *testing.T) {
			kp := mockKafkaPublisher{}
			bc := make(chan ws.EventsBatch, 2)
			m := &mockMetric{}
			m.On("Timing", "processing.latency", mock.Anything, "")
			m.On("Count", "kafka.messages.delivered", 0, "success=true")
			m.On("Count", "kafka.messages.delivered", 0, "success=false")

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
