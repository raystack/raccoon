package buffer

import (
	"errors"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockKakfaPublisher struct {
	mock.Mock
}

func (m mockKakfaPublisher) Produce(message *kafka.Message, deliveryChannel chan kafka.Event) error {
	err := m.Called(mock.Anything, mock.Anything).Error(0)
	return err
}

func TestWorker(t *testing.T) {
	t.Run("StartWorker", func(t *testing.T) {
		t.Run("Should run worker as much as poolNumbers", func(t *testing.T) {
			m := mockKakfaPublisher{}
			bc := make(chan []byte, 2)
			worker := NewWorker(1, bc, &m)
			worker.StartWorker()

			m.On("Produce", mock.Anything, mock.Anything).Return(nil).After(3 * time.Millisecond)
			bc <- []byte{}
			bc <- []byte{}
			assert.Equal(t, 1, len(bc))
		})

		t.Run("Should publish message on bufferChannel to kafka", func(t *testing.T) {
			m := mockKakfaPublisher{}
			bc := make(chan []byte, 2)
			worker := NewWorker(1, bc, &m)
			worker.StartWorker()

			m.On("Produce", mock.Anything, mock.Anything).Return(nil).Twice()
			bc <- []byte{}
			bc <- []byte{}
			time.Sleep(10 * time.Millisecond)

			m.AssertExpectations(t)
		})

		t.Run("Should retry when fail publishing to kafka", func(t *testing.T) {
			m := mockKakfaPublisher{}
			bc := make(chan []byte, 1)
			worker := NewWorker(1, bc, &m)
			worker.StartWorker()

			m.On("Produce", mock.Anything, mock.Anything).Return(errors.New("Oops")).Twice()
			m.On("Produce", mock.Anything, mock.Anything).Return(nil).Once()
			bc <- []byte{}
			time.Sleep(10 * time.Millisecond)

			m.AssertExpectations(t)
		})
	})

	t.Run("Flush", func(t *testing.T) {
		t.Run("Should block until all messages is proccessed", func(t *testing.T) {
			m := mockKakfaPublisher{}
			bc := make(chan []byte, 2)
			worker := NewWorker(1, bc, &m)
			worker.StartWorker()
			m.On("Produce", mock.Anything, mock.Anything).Return(nil).Times(3).After(3 * time.Millisecond)
			bc <- []byte{}
			bc <- []byte{}
			bc <- []byte{}
			close(bc)
			worker.Flush()
			assert.Equal(t, 0, len(bc))
		})
	})
}
