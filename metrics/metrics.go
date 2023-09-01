package metrics

import (
	"errors"

	"github.com/raystack/raccoon/config"
)

var Instrument MetricInstrument

type MetricInstrument interface {
	Increment(metricName string, labels map[string]string) error
	Count(metricName string, count int64, labels map[string]string) error
	Gauge(metricName string, value interface{}, labels map[string]string) error
	Histogram(metricName string, value int64, labels map[string]string) error
	Close()
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
		Instrument = prometheus
	}
	if config.MetricStatsd.Enabled {
		statsD, err := initStatsd()
		if err != nil {
			return err
		}
		Instrument = statsD
	}
	return nil
}

func SetVoid() {
	if config.MetricPrometheus.Enabled {
		setPrometheusVoid()
	}
	if config.MetricStatsd.Enabled {
		setStatsDVoid()
	}
}

func Close() {
	if Instrument != nil {
		Instrument.Close()
	}
}
