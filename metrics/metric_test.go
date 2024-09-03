package metrics

import (
	"os"
	"testing"
	"time"

	"github.com/raystack/raccoon/config"
	"github.com/stretchr/testify/assert"
)

func Test_Prometheus_Setup(t *testing.T) {
	config.Metric.Prometheus.Enabled = true
	config.Metric.StatsD.Enabled = false
	config.Metric.Prometheus.Path = "/metrics"
	config.Metric.Prometheus.Port = 8080
	Setup()
	prometheusInstrument, ok := instrument.(*PrometheusCollector)
	assert.True(t, ok, "prometheus collector was not initialised")
	assert.NotNil(t, prometheusInstrument)
	os.Unsetenv("METRIC_PROMETHEUS_ENABLED")
	Close()
	time.Sleep(5 * time.Second)
}

func Test_Error_On_Both_Enabled(t *testing.T) {
	config.Metric.Prometheus.Enabled = true
	config.Metric.StatsD.Enabled = true
	assert.Error(t, Setup())
	os.Setenv("METRIC_STATSD_ENABLED", "false")
	os.Setenv("METRIC_PROMETHEUS_ENABLED", "false")
	defer Close()
}

func Test_Count_Calls_Instrument_Count(t *testing.T) {
	mockInstrumentImpl := &MockInstrument{}
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
	mockInstrumentImpl := &MockInstrument{}
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
	mockInstrumentImpl := &MockInstrument{}
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
	mockInstrumentImpl := &MockInstrument{}
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
