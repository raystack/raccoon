package publisher

import (
	"github.com/stretchr/testify/mock"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
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

func (p *mockClient) Events() chan kafka.Event {
	return make(chan kafka.Event)
}
