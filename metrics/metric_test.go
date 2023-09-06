package metrics

import (
	"os"
	"testing"
	"time"

	"github.com/raystack/raccoon/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockMetricInstrument struct {
	mock.Mock
}

func (m *mockMetricInstrument) Increment(metricName string, labels map[string]string) error {
	args := m.Called(metricName, labels)
	err := args.Get(0)
	if err != nil {
		return args.Get(0).(error)
	} else {
		return nil
	}
}

func (m *mockMetricInstrument) Count(metricName string, count int64, labels map[string]string) error {
	args := m.Called(metricName, count, labels)
	err := args.Get(0)
	if err != nil {
		return args.Get(0).(error)
	} else {
		return nil
	}

}

func (m *mockMetricInstrument) Gauge(metricName string, value interface{}, labels map[string]string) error {
	args := m.Called(metricName, value, labels)
	err := args.Get(0)
	if err != nil {
		return args.Get(0).(error)
	} else {
		return nil
	}
}

func (m *mockMetricInstrument) Histogram(metricName string, value int64, labels map[string]string) error {
	args := m.Called(metricName, value, labels)
	err := args.Get(0)
	if err != nil {
		return args.Get(0).(error)
	} else {
		return nil
	}
}

func (m *mockMetricInstrument) Close() {
	m.Called()
}

func Test_Prometheus_Setup(t *testing.T) {
	config.MetricPrometheus.Enabled = true
	config.MetricStatsd.Enabled = false
	config.MetricPrometheus.Path = "/metrics"
	config.MetricPrometheus.Port = 8080
	Setup()
	prometheusInstrument, ok := instrument.(*PrometheusCollector)
	assert.True(t, ok, "prometheus collector was not initialised")
	assert.NotNil(t, prometheusInstrument)
	os.Unsetenv("METRIC_PROMETHEUS_ENABLED")
	Close()
	time.Sleep(5 * time.Second)
}

func Test_StatsDSetup(t *testing.T) {
	config.MetricPrometheus.Enabled = false
	config.MetricStatsd.Enabled = true
	config.MetricStatsd.Address = ":8125"
	config.MetricStatsd.FlushPeriodMs = 5000
	Setup()
	statsDInstrument, ok := instrument.(*Statsd)
	assert.True(t, ok, "statsd collector was not initialised")
	assert.NotNil(t, statsDInstrument)
	os.Unsetenv("METRIC_STATSD_ENABLED")
	Close()
}

func Test_Error_On_Both_Enabled(t *testing.T) {
	config.MetricPrometheus.Enabled = true
	config.MetricStatsd.Enabled = true
	assert.Error(t, Setup())
	os.Setenv("METRIC_STATSD_ENABLED", "false")
	os.Setenv("METRIC_PROMETHEUS_ENABLED", "false")
	defer Close()
}

func Test_Count_Calls_Instrument_Count(t *testing.T) {
	mockInstrumentImpl := &mockMetricInstrument{}
	instrument = mockInstrumentImpl
	metricName := "abcd"
	countValue := int64(9000)
	labels := map[string]string{"xyz": "ques", "alpha": "beta"}
	mockInstrumentImpl.On("Count", metricName, countValue, labels).Return(nil)
	err := Count(metricName, countValue, labels)
	assert.NoError(t, err)
	mockInstrumentImpl.AssertCalled(t, "Count", metricName, countValue, labels)
	mockInstrumentImpl.AssertNumberOfCalls(t, "Count", 1)
}

func Test_Gauge_Calls_Instrument_Gauge(t *testing.T) {
	mockInstrumentImpl := &mockMetricInstrument{}
	instrument = mockInstrumentImpl
	metricName := "abcd"
	countValue := int64(9000)
	labels := map[string]string{"xyz": "ques", "alpha": "beta"}
	mockInstrumentImpl.On("Gauge", metricName, countValue, labels).Return(nil)
	err := Gauge(metricName, countValue, labels)
	assert.NoError(t, err)
	mockInstrumentImpl.AssertCalled(t, "Gauge", metricName, countValue, labels)
	mockInstrumentImpl.AssertNumberOfCalls(t, "Gauge", 1)
}

func Test_Histogram_Calls_Instrument_Histogram(t *testing.T) {
	mockInstrumentImpl := &mockMetricInstrument{}
	instrument = mockInstrumentImpl
	metricName := "abcd"
	countValue := int64(9000)
	labels := map[string]string{"xyz": "ques", "alpha": "beta"}
	mockInstrumentImpl.On("Histogram", metricName, countValue, labels).Return(nil)
	err := Histogram(metricName, countValue, labels)
	assert.NoError(t, err)
	mockInstrumentImpl.AssertCalled(t, "Histogram", metricName, countValue, labels)
	mockInstrumentImpl.AssertNumberOfCalls(t, "Histogram", 1)
}

func Test_Close_Calls_Instrument_Close(t *testing.T) {
	mockInstrumentImpl := &mockMetricInstrument{}
	instrument = mockInstrumentImpl
	mockInstrumentImpl.On("Close")
	Close()
	mockInstrumentImpl.AssertCalled(t, "Close")
}

func Test_Close_Does_Not_Panic_On_Nil_Instrument(t *testing.T) {
	assert.NotPanics(t, Close)
}

func Test_ReturnsErrWhenCalledithNilInstrument(t *testing.T) {
	instrument = nil
	metricName := "abcd"
	countValue := int64(9000)
	labels := map[string]string{"xyz": "ques", "alpha": "beta"}
	err := Count(metricName, countValue, labels)
	assert.Error(t, err)
	err = Gauge(metricName, countValue, labels)
	assert.Error(t, err)
	err = Histogram(metricName, countValue, labels)
	assert.Error(t, err)
}
