package metrics

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/raystack/raccoon/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PrometheusTestSuite struct {
	suite.Suite
	instrument *PrometheusCollector
}

const prometheusPath = "/test_metrics"
const prometheusPort = "9999"

type mockCounterVec struct {
	mock.Mock
}

// With is a method of mockCounterVec that satisfies the CounterVec interface.
func (c *mockCounterVec) With(labels prometheus.Labels) prometheus.Counter {
	args := c.Called(labels)
	return args.Get(0).(prometheus.Counter)
}

// Collect is a method of mockCounterVec that satisfies the CounterVec interface.
func (c *mockCounterVec) Collect(ch chan<- prometheus.Metric) {
	c.Called(ch)
}

// Describe is a method of mockCounterVec that satisfies the CounterVec interface.
func (c *mockCounterVec) Describe(ch chan<- *prometheus.Desc) {
	c.Called(ch)
}

// mockGaugeVec is a mock struct that implements the GaugeVec interface and embeds mock.Mock.
type mockGaugeVec struct {
	mock.Mock
}

// With is a method of mockGaugeVec that satisfies the GaugeVec interface.
func (g *mockGaugeVec) With(labels prometheus.Labels) prometheus.Gauge {
	args := g.Called(labels)
	return args.Get(0).(prometheus.Gauge)
}

// Collect is a method of mockGaugeVec that satisfies the GaugeVec interface.
func (g *mockGaugeVec) Collect(ch chan<- prometheus.Metric) {
	g.Called(ch)
}

// Describe is a method of mockGaugeVec that satisfies the GaugeVec interface.
func (g *mockGaugeVec) Describe(ch chan<- *prometheus.Desc) {
	g.Called(ch)
}

// mockHistogramVec is a mock struct that implements the HistogramVec interface and embeds mock.Mock.
type mockHistogramVec struct {
	mock.Mock
}

// With is a method of mockHistogramVec that satisfies the HistogramVec interface.
func (h *mockHistogramVec) With(labels prometheus.Labels) prometheus.Observer {
	args := h.Called(labels)
	return args.Get(0).(prometheus.Observer)
}

// Collect is a method of mockHistogramVec that satisfies the HistogramVec interface.
func (h *mockHistogramVec) Collect(ch chan<- prometheus.Metric) {
	h.Called(ch)
}

// Describe is a method of mockHistogramVec that satisfies the HistogramVec interface.
func (h *mockHistogramVec) Describe(ch chan<- *prometheus.Desc) {
	h.Called(ch)
}

type mockCounter struct {
	mock.Mock
	prometheus.Counter
}

func (m *mockCounter) Add(f float64) {
	m.Called(f)
}

func (m *mockCounter) Inc() {
	m.Called()
}

type mockGauge struct {
	mock.Mock
	prometheus.Gauge
}

func (m *mockGauge) Set(f float64) {
	m.Called(f)
}

type mockObserver struct {
	mock.Mock
}

func (m *mockObserver) Observe(f float64) {
	m.Called(f)
}

func (promSuite *PrometheusTestSuite) Test_Prometheus_Collector_Metrics_Initialised() {
	// NOTE(turtledev): what are we even testing here?
	numCounters := 22
	numGauge := 15
	numHistogram := 10
	var err error
	promSuite.instrument, err = initPrometheusCollector()
	assert.NoError(promSuite.T(), err, "error while initialising prometheus collector")
	assert.Equal(promSuite.T(), numCounters, len(promSuite.instrument.counters), "number of prometheus counters do not match expected count")
	assert.Equal(promSuite.T(), numGauge, len(promSuite.instrument.gauges), "number of prometheus gauges do not match expected count")
	assert.Equal(promSuite.T(), numHistogram, len(promSuite.instrument.histogram), "number of prometheus histogram do not match expected count")
}

func (promSuite *PrometheusTestSuite) Test_PrometheusCollector_Metrics_Registered() {
	var err error
	promSuite.instrument, err = initPrometheusCollector()
	assert.NoError(promSuite.T(), err)
	for metricName, collector := range promSuite.instrument.counters {
		assert.True(promSuite.T(), promSuite.instrument.registry.Unregister(collector), fmt.Sprintf("%s counter not registered", metricName))
	}
	for metricName, collector := range promSuite.instrument.gauges {
		assert.True(promSuite.T(), promSuite.instrument.registry.Unregister(collector), fmt.Sprintf("%s gauge not registered", metricName))
	}
	for metricName, collector := range promSuite.instrument.histogram {
		assert.True(promSuite.T(), promSuite.instrument.registry.Unregister(collector), fmt.Sprintf("%s counter not registered", metricName))
	}
}

func (promSuite *PrometheusTestSuite) Test_Prometheus_Collector_MetricServer_Initialised() {
	var err error
	promSuite.instrument, err = initPrometheusCollector()
	assert.NoError(promSuite.T(), err)
	client := http.Client{Timeout: 5 * time.Second}
	time.Sleep(2 * time.Second)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://:%s%s", prometheusPort, prometheusPath), nil)
	assert.NoError(promSuite.T(), err)
	res, err := client.Do(req)
	assert.NoError(promSuite.T(), err)
	defer io.Copy(io.Discard, res.Body)
	defer res.Body.Close()
	assert.Equal(promSuite.T(), http.StatusOK, res.StatusCode)
	bodyLength, err := io.Copy(io.Discard, res.Body)
	assert.NoError(promSuite.T(), err)
	assert.NotZero(promSuite.T(), bodyLength)
}

func (promSuite *PrometheusTestSuite) Test_Prometheus_Counter_Working() {
	sampleCounterMetric1 := "kafka_unknown_topic_failure_total"
	sampleCounterMetric2 := "batches_read_total"
	mockCounterVec1 := mockCounterVec{}
	mockCounterVec2 := mockCounterVec{}
	mockCounter1 := mockCounter{}
	mockCounter2 := mockCounter{}
	callArg1 := int64(35)
	labels1 := map[string]string{"topic": "test", "event_type": "abc", "conn_group": "--default--"}
	promLabels1 := prometheus.Labels{"topic": "test", "event_type": "abc", "conn_group": "--default--"}
	labels2 := map[string]string{"status": "success", "reason": "unknown", "conn_group": "abc"}
	promLabels2 := prometheus.Labels{"status": "success", "reason": "unknown", "conn_group": "abc"}
	var err error
	promSuite.instrument, err = initPrometheusCollector()
	assert.NoError(promSuite.T(), err)
	promSuite.instrument.counters[sampleCounterMetric1] = &mockCounterVec1
	promSuite.instrument.counters[sampleCounterMetric2] = &mockCounterVec2
	mockCounterVec1.On("With", promLabels1).Return(&mockCounter1)
	mockCounter1.On("Add", float64(callArg1))
	mockCounterVec2.On("With", promLabels2).Return(&mockCounter2)
	mockCounter2.On("Inc")
	err = promSuite.instrument.Count(sampleCounterMetric1, callArg1, labels1)
	assert.NoError(promSuite.T(), err)
	err = promSuite.instrument.Increment(sampleCounterMetric2, labels2)
	assert.NoError(promSuite.T(), err)
	mockCounterVec1.AssertCalled(promSuite.T(), "With", promLabels1)
	mockCounterVec2.AssertCalled(promSuite.T(), "With", promLabels2)
	mockCounterVec1.AssertNumberOfCalls(promSuite.T(), "With", 1)
	mockCounterVec2.AssertNumberOfCalls(promSuite.T(), "With", 1)
	mockCounter1.AssertNumberOfCalls(promSuite.T(), "Add", 1)
	mockCounter2.AssertNumberOfCalls(promSuite.T(), "Inc", 1)
	mockCounter1.AssertCalled(promSuite.T(), "Add", float64(callArg1))
	mockCounter2.AssertCalled(promSuite.T(), "Inc")
}

func (promSuite *PrometheusTestSuite) Test_Prometheus_Gauge_Working() {
	sampleGaugeMetric := "server_go_routines_count_current"
	mockGaugeVec := mockGaugeVec{}
	mockGauge := mockGauge{}
	callArg := int64(35)
	labels := map[string]string{"topic": "test", "event_type": "abc"}
	promLabels := prometheus.Labels{"topic": "test", "event_type": "abc"}
	var err error
	promSuite.instrument, err = initPrometheusCollector()
	assert.NoError(promSuite.T(), err)
	promSuite.instrument.gauges[sampleGaugeMetric] = &mockGaugeVec
	mockGaugeVec.On("With", promLabels).Return(&mockGauge)
	mockGauge.On("Set", float64(callArg))
	err = promSuite.instrument.Gauge(sampleGaugeMetric, callArg, labels)
	assert.NoError(promSuite.T(), err)
	mockGaugeVec.AssertCalled(promSuite.T(), "With", promLabels)
	mockGaugeVec.AssertNumberOfCalls(promSuite.T(), "With", 1)
	mockGauge.AssertNumberOfCalls(promSuite.T(), "Set", 1)
	mockGauge.AssertCalled(promSuite.T(), "Set", float64(callArg))
}

func (promSuite *PrometheusTestSuite) Test_Prometheus_Histogram_Working() {
	sampleHistogramMetric := "server_processing_latency_milliseconds"
	mockHistogramVec := mockHistogramVec{}
	mockObserver := mockObserver{}
	callArg := int64(35)
	labels := map[string]string{"conn_group": "--default--"}
	promLabels := prometheus.Labels{"conn_group": "--default--"}
	var err error
	promSuite.instrument, err = initPrometheusCollector()
	assert.NoError(promSuite.T(), err)
	promSuite.instrument.histogram[sampleHistogramMetric] = &mockHistogramVec
	mockHistogramVec.On("With", promLabels).Return(&mockObserver)
	mockObserver.On("Observe", float64(callArg))
	err = promSuite.instrument.Histogram(sampleHistogramMetric, callArg, labels)
	assert.NoError(promSuite.T(), err)
	mockHistogramVec.AssertCalled(promSuite.T(), "With", promLabels)
	mockHistogramVec.AssertNumberOfCalls(promSuite.T(), "With", 1)
	mockObserver.AssertCalled(promSuite.T(), "Observe", float64(callArg))
	mockObserver.AssertNumberOfCalls(promSuite.T(), "Observe", 1)
}

func (promSuite *PrometheusTestSuite) Test_Prometheus_Gauge_Error_On_Invalid_Input() {
	sampleGaugeMetric := "server_go_routines_count_current"
	mockGaugeVec := mockGaugeVec{}
	callArg := "abc"
	labels := map[string]string{"topic": "test", "event_type": "abc"}
	var err error
	promSuite.instrument, err = initPrometheusCollector()
	assert.NoError(promSuite.T(), err)
	promSuite.instrument.gauges[sampleGaugeMetric] = &mockGaugeVec
	err = promSuite.instrument.Gauge(sampleGaugeMetric, callArg, labels)
	assert.Error(promSuite.T(), err)
	mockGaugeVec.AssertNotCalled(promSuite.T(), "With")
}

func TestPrometheusSuite(t *testing.T) {
	suite.Run(t, new(PrometheusTestSuite))
}

func (suite *PrometheusTestSuite) SetupTest() {
	var err error
	config.Metric.Prometheus.Enabled = true
	config.Metric.Prometheus.Path = prometheusPath
	config.Metric.Prometheus.Port, err = strconv.Atoi(prometheusPort)
	assert.NoError(suite.T(), err)
}

func (promSuite *PrometheusTestSuite) TearDownTest() {
	if promSuite.instrument != nil {
		promSuite.instrument.Close()
	}
}
