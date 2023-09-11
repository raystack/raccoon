package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/logger"
	"github.com/spf13/cast"
)

type CounterVec interface {
	With(labels prometheus.Labels) prometheus.Counter
	Collect(chan<- prometheus.Metric)
	Describe(chan<- *prometheus.Desc)
}

type GaugeVec interface {
	With(labels prometheus.Labels) prometheus.Gauge
	Collect(chan<- prometheus.Metric)
	Describe(chan<- *prometheus.Desc)
}

type HistogramVec interface {
	With(labels prometheus.Labels) prometheus.Observer
	Collect(chan<- prometheus.Metric)
	Describe(chan<- *prometheus.Desc)
}

type PrometheusCollector struct {
	registry  *prometheus.Registry
	counters  map[string]CounterVec
	gauges    map[string]GaugeVec
	histogram map[string]HistogramVec
	server    *http.Server
}

func initPrometheusCollector() (*PrometheusCollector, error) {
	serveMux := &http.ServeMux{}
	server := &http.Server{Addr: fmt.Sprintf(":%d", config.MetricPrometheus.Port), Handler: serveMux}
	p := &PrometheusCollector{
		counters:  getCounterMap(),
		gauges:    getGaugeMap(),
		histogram: getHistogramMap(),
		registry:  prometheus.NewRegistry(),
		server:    server,
	}
	err := p.Register()
	if err != nil {
		return nil, err
	}
	serveMux.Handle(config.MetricPrometheus.Path, promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{Registry: p.registry}))
	go server.ListenAndServe()
	return p, nil

}

func (p *PrometheusCollector) Count(metricName string, count int64, labels map[string]string) error {
	counter, ok := p.counters[metricName]
	if !ok {
		return fmt.Errorf("invalid counter metric %s", metricName)
	}
	counter.With(labels).Add(float64(count))
	return nil
}

func (p *PrometheusCollector) Increment(metricName string, labels map[string]string) error {
	counter, ok := p.counters[metricName]
	if !ok {
		return fmt.Errorf("invalid counter metric %s", metricName)
	}
	counter.With(labels).Inc()
	return nil
}

func (p *PrometheusCollector) Gauge(metricName string, value interface{}, labels map[string]string) error {
	gauge, ok := p.gauges[metricName]
	if !ok {
		return fmt.Errorf("invalid gauge metric %s", metricName)
	}
	floatVal, err := cast.ToFloat64E(value)
	if err != nil {
		return err
	}
	gauge.With(labels).Set(floatVal)
	return nil
}

func (p *PrometheusCollector) Histogram(metricName string, value int64, labels map[string]string) error {
	histogram, ok := p.histogram[metricName]
	if !ok {
		return fmt.Errorf("invalid histogram metric %s", metricName)
	}
	floatVal, err := cast.ToFloat64E(value)
	if err != nil {
		return err
	}
	histogram.With(labels).Observe(floatVal)
	return nil
}

func (p *PrometheusCollector) Register() error {
	for _, counter := range p.counters {
		err := p.registry.Register(counter)
		if err != nil {
			return err
		}
	}
	for _, gauge := range p.gauges {
		err := p.registry.Register(gauge)
		if err != nil {
			return err
		}
	}
	for _, histogram := range p.histogram {
		err := p.registry.Register(histogram)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PrometheusCollector) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	logger.Info("shutting down prometheus metric server")
	err := p.server.Shutdown(ctx)
	if err != nil {
		logger.Warn("error shutting down logger")
	}
	defer cancel()
}

func getCounterMap() map[string]CounterVec {
	counters := make(map[string]CounterVec)
	counters["kafka_messages_delivered_total"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "kafka_messages_delivered_total",
		Help: "Number of delivered events to Kafka"}, []string{"success", "conn_group", "event_type"})
	counters["kafka_messages_undelivered_total"] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "kafka_messages_undelivered_total",
		Help: "Number of delivered events to Kafka which failed while reading delivery report"}, []string{"success", "conn_group", "event_type"})
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

func getGaugeMap() map[string]GaugeVec {
	gauges := make(map[string]GaugeVec)
	gauges["server_go_routines_count_current"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "server_go_routines_count_current",
		Help: "Number of goroutine spawn in a single flush"}, []string{})
	gauges["server_mem_heap_alloc_bytes_current"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "server_mem_heap_alloc_bytes_current",
		Help: "Bytes of allocated heap objects"}, []string{})
	gauges["server_mem_heap_inuse_bytes_current"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "server_mem_heap_inuse_bytes_current",
		Help: "HeapInuse is bytes in in-use spans"}, []string{})
	gauges["server_mem_heap_objects_total_current"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "server_mem_heap_objects_total_current",
		Help: "Number of allocated heap objects"}, []string{})
	gauges["server_mem_stack_inuse_bytes_current"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "server_mem_stack_inuse_bytes_current",
		Help: "Bytes in stack spans"}, []string{})
	gauges["server_mem_gc_triggered_current"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "server_mem_gc_triggered_current",
		Help: "The time the last garbage collection finished in Unix timestamp"}, []string{})
	gauges["server_mem_gc_pauseNs_current"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "server_mem_gc_pauseNs_current",
		Help: "Circular buffer of recent GC stop-the-world in Unix timestamp"}, []string{})
	gauges["server_mem_gc_count_current"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "server_mem_gc_count_current",
		Help: "The number of completed GC cycle"}, []string{})
	gauges["server_mem_gc_pauseTotalNs_current"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "server_mem_gc_pauseTotalNs_current",
		Help: "The cumulative nanoseconds in GC stop-the-world pauses since the program started"}, []string{})
	gauges["kafka_tx_messages_total"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kafka_tx_messages_total",
		Help: "Total number of messages transmitted produced to Kafka brokers."}, []string{})
	gauges["kafka_tx_messages_bytes_total"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kafka_tx_messages_bytes_total",
		Help: "Total number of message bytes (including framing, such as per-Message framing and MessageSet/batch framing) transmitted to Kafka brokers"}, []string{})
	gauges["kafka_brokers_tx_total"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kafka_brokers_tx_total",
		Help: "Total number of requests sent to Kafka brokers"}, []string{"broker"})
	gauges["kafka_brokers_tx_bytes_total"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kafka_brokers_tx_bytes_total",
		Help: "Total number of bytes transmitted to Kafka brokers"}, []string{"broker"})
	gauges["kafka_brokers_rtt_average_milliseconds"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kafka_brokers_rtt_average_milliseconds",
		Help: "Broker latency / round-trip time in microseconds"}, []string{"broker"})
	gauges["connections_count_current"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "connections_count_current",
		Help: "Number of alive connections"}, []string{"conn_group"})
	return gauges
}

func getHistogramMap() map[string]HistogramVec {
	histogram := make(map[string]HistogramVec)
	histogram["ack_event_rtt_ms"] = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ack_event_rtt_ms",
		Help: "Time taken from ack function called by kafka producer to processed by the ack handler.",
		Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000},
	}, []string{})
	histogram["event_rtt_ms"] = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "event_rtt_ms",
		Help: "Time taken from event is consumed from the queue to be acked by the ack handler.",
		Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000},
	}, []string{})
	histogram["user_session_duration_milliseconds"] = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "user_session_duration_milliseconds",
		Help:    "Duration of alive connection per session per connection",
		Buckets: []float64{1000, 10000, 100000, 600000, 3600000},
	}, []string{"conn_group"})
	histogram["batch_idle_in_channel_milliseconds"] = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "batch_idle_in_channel_milliseconds",
		Help:    "Duration from when the request is received to when the request is processed. High value of this metric indicates the publisher is slow.",
		Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000},
	}, []string{"worker"})
	histogram["kafka_producebulk_tt_ms"] = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "kafka_producebulk_tt_ms",
		Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000},
	}, []string{})
	histogram["event_processing_duration_milliseconds"] = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "event_processing_duration_milliseconds",
		Help:    "Duration from the time request is sent to the time events are published. This metric is calculated per event by following formula (PublishedTime - SentTime)/CountEvents",
		Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000},
	}, []string{"conn_group"})
	histogram["worker_processing_duration_milliseconds"] = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "worker_processing_duration_milliseconds",
		Help:    "Duration from the time request is processed to the time events are published. This metric is calculated per event by following formula (PublishedTime - ProcessedTime)/CountEvents",
		Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000},
	}, []string{"worker"})
	histogram["server_processing_latency_milliseconds"] = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "server_processing_latency_milliseconds",
		Help:    "Duration from the time request is receieved to the time events are published. This metric is calculated per event by following formula`(PublishedTime - ReceievedTime)/CountEvents`",
		Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000},
	}, []string{"conn_group"})
	return histogram
}
