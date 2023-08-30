package metrics

type MetricInstrument interface {
	Increment(metricName string, labels map[string]string)
	Count(metricName string, labels map[string]string, count int64)
	Gauge(metricName string, labels map[string]string, value float64)
	Histogram(metricName string, labels map[string]string, value float64)
}
