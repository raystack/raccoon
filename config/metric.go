package config

import (
	"time"

	"github.com/raystack/raccoon/config/util"

	"github.com/spf13/viper"
)

var MetricStatsd metricStatsdCfg
var MetricPrometheus metricPrometheusCfg
var MetricInfo metricInfoCfg

type metricStatsdCfg struct {
	Enabled       bool
	Address       string
	FlushPeriodMs time.Duration
}

type metricPrometheusCfg struct {
	Enabled bool
	Port    int
	Path    string
}

type metricInfoCfg struct {
	RuntimeStatsRecordInterval time.Duration
}

func metricStatsdConfigLoader() {
	viper.SetDefault("METRIC_STATSD_ENABLED", false)
	viper.SetDefault("METRIC_STATSD_ADDRESS", ":8125")
	viper.SetDefault("METRIC_STATSD_FLUSH_PERIOD_MS", 10000)
	MetricStatsd = metricStatsdCfg{
		Enabled:       util.MustGetBool("METRIC_STATSD_ENABLED"),
		Address:       util.MustGetString("METRIC_STATSD_ADDRESS"),
		FlushPeriodMs: util.MustGetDuration("METRIC_STATSD_FLUSH_PERIOD_MS", time.Millisecond),
	}
}

func metricPrometheusConfigLoader() {
	viper.SetDefault("METRIC_PROMETHEUS_ENABLED", false)
	viper.SetDefault("METRIC_PROMETHEUS_PORT", 9090)
	viper.SetDefault("METRIC_PROMETHEUS_PATH", "/metrics")
	MetricPrometheus = metricPrometheusCfg{
		Enabled: util.MustGetBool("METRIC_PROMETHEUS_ENABLED"),
		Port:    util.MustGetInt("METRIC_PROMETHEUS_PORT"),
		Path:    util.MustGetString("METRIC_PROMETHEUS_PATH"),
	}
}

func metricCommonConfigLoader() {
	viper.SetDefault("METRIC_RUNTIME_STATS_RECORD_INTERVAL_MS", 10000)
	MetricInfo = metricInfoCfg{
		RuntimeStatsRecordInterval: util.MustGetDuration("METRIC_RUNTIME_STATS_RECORD_INTERVAL_MS", time.Millisecond),
	}
}
