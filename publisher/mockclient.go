package publisher

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/mock"
)

type mockClient struct {
	mock.Mock
}

func (p *mockClient) Produce(m *kafka.Message, eventsChan chan kafka.Event) error {
	args := p.Called(m, eventsChan)
	return args.Error(0)
}

func (p *mockClient) Close() {
	p.Called()
}

func (p *mockClient) Flush(config int) int {
	args := p.Called(config)
	return args.Int(0)
}
