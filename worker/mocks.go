package worker

import (
	"github.com/stretchr/testify/mock"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"source.golabs.io/mobile/clickstream-go-proto/gojek/clickstream/de"
)

type mockKafkaPublisher struct {
	mock.Mock
}

func (m *mockKafkaPublisher) ProduceBulk(events []*de.Event, deliveryChannel chan kafka.Event) error {
	mock := m.Called(events, deliveryChannel)
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
