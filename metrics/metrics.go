package metrics

type MetricInstrument interface {
	Increment(metricName string, labels map[string]string) error
	Count(metricName string, labels map[string]string, count int64) error
	Gauge(metricName string, labels map[string]string, value float64) error
	Histogram(metricName string, labels map[string]string, value float64) error
}
