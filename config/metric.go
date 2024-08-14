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
	Enabled     bool          `mapstructure:"METRIC_STATSD_ENABLED" cmdx:"metric.statsd.enabled" default:"false" `
	Address     string        `mapstructure:"METRIC_STATSD_ADDRESS" cmdx:"metric.statsd.address" default:":8125"`
	FlushPeriod time.Duration `mapstructure:"METRIC_STATSD_FLUSH_PERIOD_MS" cmdx:"metric.statsd.flush.period.ms" default:"10000"`
}

type metricPrometheusCfg struct {
	Enabled bool   `mapstructure:"METRIC_PROMETHEUS_ENABLED" cmdx:"metric.prometheus.enabled" default:"false"`
	Port    int    `mapstructure:"METRIC_PROMETHEUS_PORT" cmdx:"metric.prometheus.port" default:"9090"`
	Path    string `mapstructure:"METRIC_PROMETHEUS_PATH" cmdx:"metric.prometheus.path" default:"/metrics"`
}

type metricInfoCfg struct {
	RuntimeStatsRecordInterval time.Duration `mapstructure:"METRIC_RUNTIME_STATS_RECORD_INTERVAL_MS" cmdx:"metric.runtime.stats.record.interval.ms" default:"10000"`
}

func metricStatsdConfigLoader() {
	viper.SetDefault("METRIC_STATSD_ENABLED", false)
	viper.SetDefault("METRIC_STATSD_ADDRESS", ":8125")
	viper.SetDefault("METRIC_STATSD_FLUSH_PERIOD_MS", 10000)
	MetricStatsd = metricStatsdCfg{
		Enabled:     util.MustGetBool("METRIC_STATSD_ENABLED"),
		Address:     util.MustGetString("METRIC_STATSD_ADDRESS"),
		FlushPeriod: util.MustGetDuration("METRIC_STATSD_FLUSH_PERIOD_MS", time.Millisecond),
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
