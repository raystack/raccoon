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
		logger.Warn("instrumentation is not set for logging")
		return errors.New("instrumentation is not set for logging")
	}
	return instrument.Increment(metricName, labels)
}

func Count(metricName string, count int64, labels map[string]string) error {
	if instrument == nil {
		logger.Warn("instrumentation is not set for logging")
		return errors.New("instrumentation is not set for logging")
	}
	return instrument.Count(metricName, count, labels)
}

func Gauge(metricName string, value interface{}, labels map[string]string) error {
	if instrument == nil {
		logger.Warn("instrumentation is not set for logging")
		return errors.New("instrumentation is not set for logging")
	}
	return instrument.Gauge(metricName, value, labels)
}

func Histogram(metricName string, value int64, labels map[string]string) error {
	if instrument == nil {
		logger.Warn("instrumentation is not set for logging")
		return errors.New("instrumentation is not set for logging")
	}
	return instrument.Histogram(metricName, value, labels)
}

func Setup() error {

	if config.MetricPrometheus.Enabled && config.MetricStatsd.Enabled {
		return errors.New("only one instrumentation tool can be enabled")
	}

	if config.MetricPrometheus.Enabled {
		prometheus, err := initPrometheusCollector()
		if err != nil {
			return err
		}
		instrument = prometheus
	}
	if config.MetricStatsd.Enabled {
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

func Close() {
	if instrument != nil {
		instrument.Close()
	}
}
