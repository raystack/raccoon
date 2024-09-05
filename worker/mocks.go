package worker

import (
	pb "github.com/raystack/raccoon/proto"
	mock "github.com/stretchr/testify/mock"
)

// KafkaProducer is an autogenerated mock type for the KafkaProducer type
type mockKafkaPublisher struct {
	mock.Mock
}

// ProduceBulk provides a mock function with given fields: events, deliveryChannel
func (m *mockKafkaPublisher) ProduceBulk(events []*pb.Event, connGroup string) error {
	mock := m.Called(events, connGroup)
	return mock.Error(0)
}

func (m *mockKafkaPublisher) Name() string {
	return m.Called().String(0)
}

type mockAck struct {
	mock.Mock
}

func (m *mockAck) Ack(err error) {
	m.Called(err)
}
