package publisher

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/mock"
)

type MockKafkaProducer struct {
	mock.Mock
}

func (p *MockKafkaProducer) Produce(m *kafka.Message, eventsChan chan kafka.Event) error {
	args := p.Called(m, eventsChan)
	return args.Error(0)
}

func (p *MockKafkaProducer) Close() {
	p.Called()
}

