package metrics

import "github.com/stretchr/testify/mock"

type MockInstrument struct {
	mock.Mock
}

func (m *MockInstrument) Increment(metricName string, labels map[string]string) error {
	args := m.Called(metricName, labels)
	err := args.Get(0)
	if err != nil {
		return args.Get(0).(error)
	} else {
		return nil
	}
}

func (m *MockInstrument) Count(metricName string, count int64, labels map[string]string) error {
	args := m.Called(metricName, count, labels)
	err := args.Get(0)
	if err != nil {
		return args.Get(0).(error)
	} else {
		return nil
	}
}

func (m *MockInstrument) Gauge(metricName string, value interface{}, labels map[string]string) error {
	args := m.Called(metricName, value, labels)
	err := args.Get(0)
	if err != nil {
		return args.Get(0).(error)
	} else {
		return nil
	}
}

func (m *MockInstrument) Histogram(metricName string, value int64, labels map[string]string) error {
	args := m.Called(metricName, value, labels)
	err := args.Get(0)
	if err != nil {
		return args.Get(0).(error)
	} else {
		return nil
	}
}

func (m *MockInstrument) Close() {
	m.Called()
}
