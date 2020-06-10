package publisher

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/mock"
)

type MockProducer struct {
	mock.Mock
}

func (p *MockProducer) Produce(msg *kafka.Message) error {
	args := p.Called(msg)
	return args.Error(0)
}

