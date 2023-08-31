package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusCollector struct {
	counters  map[string]*prometheus.CounterVec
	gauges    map[string]prometheus.Gauge
	histogram map[string]prometheus.Histogram
}

func initPrometheusCollector() *PrometheusCollector {
	return &PrometheusCollector{
		counters:  getCounterMap(),
		gauges:    getGaugeMap(),
		histogram: getHistogramMap(),
	}
}

func (p *PrometheusCollector) Count(metricName string, labels map[string]string, count int64) error {
	counter, ok := p.counters[metricName]
	if !ok {
		return fmt.Errorf("invalid counter metric %s", metricName)
	}
	counter.With(labels).Add(float64(count))
	return nil
}

func (p *PrometheusCollector) Increment(metricName string, labels map[string]string) error {
	return nil
}

func (p *PrometheusCollector) Gauge(metricName string, labels map[string]string, value float64) error {
	return nil
}

func (p *PrometheusCollector) Histogram(metricName string, labels map[string]string, value float64) error {
	return nil
}

func getCounterMap() map[string]*prometheus.CounterVec {
	counters := make(map[string]*prometheus.CounterVec)
	counters["kafka_messages_delivered_total"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "kafka_messages_delivered_total",
		Help: "Number of delivered events to Kafka"}, []string{"success", "conn_group", "event_type"})
	counters["events_rx_bytes_total"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "events_rx_bytes_total",
		Help: "Total byte receieved in requests"}, []string{"conn_group", "event_type"})
	counters["events_rx_total"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "events_rx_total",
		Help: "Number of events received in requests"}, []string{"conn_group", "event_type"})
	counters["kafka_unknown_topic_failure_total"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "kafka_unknown_topic_failure_total",
		Help: "Number of delivery failure caused by topic does not exist in kafka."}, []string{"topic", "event_type"})
	counters["batches_read_total"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "batches_read_total",
		Help: "Request count"}, []string{"status", "reason", "conn_group"})
	counters["events_duplicate_total"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "events_duplicate_total",
		Help: "Total Number of Duplicate events recieved by the server"}, []string{"conn_group", "reason"})
	counters["server_ping_failure_total"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "server_ping_failure_total",
		Help: "Total ping that server fails to send"}, []string{"conn_group"})
	counters["conn_close_err_count"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "conn_close_err_count",
		Help: "Total Number of Connection Errors When Trying To Close"}, []string{})
	counters["user_connection_failure_total"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "user_connection_failure_total",
		Help: "Number of fail connections established to the server"}, []string{"conn_group", "reason"})
	counters["user_connection_success_total"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "user_connection_success_total",
		Help: "Number of successful connections established to the server"}, []string{"conn_group"})
	counters["server_pong_failure_total"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "server_pong_failure_total",
		Help: "Total pong that server fails to send"}, []string{"conn_group"})
	return counters
}

func getGaugeMap() map[string]prometheus.Gauge {
	gauges := make(map[string]prometheus.Gauge)
	return gauges
}

func getHistogramMap() map[string]prometheus.Histogram {
	histogram := make(map[string]prometheus.Histogram)
	return histogram
}
