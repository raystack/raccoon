package metrics

import (
	"errors"

	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/logger"
)

var instrument MetricInstrument

type MetricInstrument interface {
	Increment(metricName string, labels map[string]string) error
	Count(metricName string, count int64, labels map[string]string) error
	Gauge(metricName string, value interface{}, labels map[string]string) error
	Histogram(metricName string, value int64, labels map[string]string) error
	Close()
}

type voidInstrument struct{}

func (v voidInstrument) Increment(metricName string, labels map[string]string) error {
	return nil
}

func (v voidInstrument) Count(metricName string, count int64, labels map[string]string) error {
	return nil
}

func (v voidInstrument) Gauge(metricName string, value interface{}, labels map[string]string) error {
	return nil
}

func (v voidInstrument) Histogram(metricName string, value int64, labels map[string]string) error {
	return nil
}

func (v voidInstrument) Close() {}

func Increment(metricName string, labels map[string]string) error {
	if instrument == nil {
		logger.Warn("instrumentation is not set for recording metrics")
		return errors.New("instrumentation is not set for recording metrics")
	}
	err := instrument.Increment(metricName, labels)
	if err != nil {
		logger.Warn(err)
	}
	return err
}

func Count(metricName string, count int64, labels map[string]string) error {
	if instrument == nil {
		logger.Warn("instrumentation is not set for recording metrics")
		return errors.New("instrumentation is not set for recording metrics")
	}
	err := instrument.Count(metricName, count, labels)
	if err != nil {
		logger.Warn(err)
	}
	return err
}

func Gauge(metricName string, value interface{}, labels map[string]string) error {
	if instrument == nil {
		logger.Warn("instrumentation is not set for recording metrics")
		return errors.New("instrumentation is not set for recording metrics")
	}
	err := instrument.Gauge(metricName, value, labels)
	if err != nil {
		logger.Warn(err)
	}
	return err
}

func Histogram(metricName string, value int64, labels map[string]string) error {
	if instrument == nil {
		logger.Warn("instrumentation is not set for recording metrics")
		return errors.New("instrumentation is not set for recording metrics")
	}
	err := instrument.Histogram(metricName, value, labels)
	if err != nil {
		logger.Warn(err)
	}
	return err
}

func Setup() error {
	if config.Metric.Prometheus.Enabled && config.Metric.StatsD.Enabled {
		return errors.New("only one instrumentation tool can be enabled")
	}

	if config.Metric.Prometheus.Enabled {
		prometheus, err := initPrometheusCollector()
		if err != nil {
			return err
		}
		instrument = prometheus
	}
	if config.Metric.StatsD.Enabled {
		statsD, err := initStatsd()
		if err != nil {
			return err
		}
		instrument = statsD
	}
	return nil
}

func SetVoid() {
	instrument = voidInstrument{}
}

// Instrument returns the configured MetricInstrument
// should be called once Setup() is done.
func Instrument() MetricInstrument {
	return instrument
}

func Close() {
	if instrument != nil {
		instrument.Close()
	}
}
