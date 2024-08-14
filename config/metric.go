package config

var MetricStatsd metricStatsdCfg
var MetricPrometheus metricPrometheusCfg
var MetricInfo metricInfoCfg

type metricStatsdCfg struct {
	Enabled       bool   `mapstructure:"METRIC_STATSD_ENABLED" cmdx:"metric.statsd.enabled" default:"false" `
	Address       string `mapstructure:"METRIC_STATSD_ADDRESS" cmdx:"metric.statsd.address" default:":8125"`
	FlushPeriodMS int64  `mapstructure:"METRIC_STATSD_FLUSH_PERIOD_MS" cmdx:"metric.statsd.flush.period.ms" default:"10000"`
}

type metricPrometheusCfg struct {
	Enabled bool   `mapstructure:"METRIC_PROMETHEUS_ENABLED" cmdx:"metric.prometheus.enabled" default:"false"`
	Port    int    `mapstructure:"METRIC_PROMETHEUS_PORT" cmdx:"metric.prometheus.port" default:"9090"`
	Path    string `mapstructure:"METRIC_PROMETHEUS_PATH" cmdx:"metric.prometheus.path" default:"/metrics"`
}

type metricInfoCfg struct {
	RuntimeStatsRecordIntervalMS int64 `mapstructure:"METRIC_RUNTIME_STATS_RECORD_INTERVAL_MS" cmdx:"metric.runtime.stats.record.interval.ms" default:"10000"`
}
