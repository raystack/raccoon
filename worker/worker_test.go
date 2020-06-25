package worker

import (
	"testing"
	"time"

	"source.golabs.io/mobile/clickstream-go-proto/gojek/clickstream/de"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type mockKakfaPublisher struct {
	mock.Mock
}

func (m *mockKakfaPublisher) ProduceBulk(message [][]byte, deliveryChannel chan kafka.Event) error {
	mock := m.Called(mock.Anything, mock.Anything)
	return mock.Error(0)
}

func TestWorker(t *testing.T) {
	t.Run("StartWorkers", func(t *testing.T) {
		t.Run("Should publish messages on bufferChannel to kafka", func(t *testing.T) {
			m := mockKakfaPublisher{}
			bc := make(chan []*de.CSEventMessage, 2)
			worker := CreateWorkerPool(1, bc, 100, &m)
			worker.StartWorkers()

			m.On("ProduceBulk", mock.Anything, mock.Anything).Return(nil).Twice()
			bc <- []*de.CSEventMessage{}
			bc <- []*de.CSEventMessage{}
			time.Sleep(10 * time.Millisecond)

			m.AssertExpectations(t)
		})
	})

	t.Run("Flush", func(t *testing.T) {
		t.Run("Should block until all messages is processed", func(t *testing.T) {
			m := mockKakfaPublisher{}
			bc := make(chan []*de.CSEventMessage, 2)
			worker := CreateWorkerPool(1, bc, 100, &m)
			worker.StartWorkers()
			m.On("ProduceBulk", mock.Anything, mock.Anything).Return(nil).Times(3).After(3 * time.Millisecond)
			bc <- []*de.CSEventMessage{}
			bc <- []*de.CSEventMessage{}
			bc <- []*de.CSEventMessage{}
			close(bc)
			worker.Flush()
			assert.Equal(t, 0, len(bc))
		})
	})
}
