package worker

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"testing"
)

type mockKakfaPublisher struct {
	mock.Mock
}

func (m *mockKakfaPublisher) ProduceBulk(message [][]byte, deliveryChannel chan kafka.Event) error {
	mock := m.Called(mock.Anything, mock.Anything)
	return mock.Error(0)
}

type mockMetric struct {
	mock.Mock
}

func (m *mockMetric) Count(bucket string, val int, tags string) {
	m.Called(bucket, val, tags)
}

func (m *mockMetric) Timing(bucket string, t int64, tags string) {
	m.Called(bucket, t, tags)
}

func mockInstrumentation(t *testing.T, xTotal, xProcessed, xErr int, xSenttime timestamp.Timestamp) func(total int, processed int, err int, sentTime timestamp.Timestamp) {
	return func(total int, processed int, err int, sentTime timestamp.Timestamp) {
		assert.Equal(t, xTotal, total)
		assert.Equal(t, xProcessed, processed)
		assert.Equal(t, xErr, err)
		assert.Equal(t, xSenttime, sentTime)
	}
}
