package config

var Metric metric

type metricStatsdCfg struct {
	Enabled       bool   `mapstructure:"enabled" cmdx:"metric.statsd.enabled" default:"false" `
	Address       string `mapstructure:"address" cmdx:"metric.statsd.address" default:":8125"`
	FlushPeriodMS int64  `mapstructure:"flush_period_ms" cmdx:"metric.statsd.flush.period.ms" default:"10000"`
}

type metricPrometheusCfg struct {
	Enabled bool   `mapstructure:"enabled" cmdx:"metric.prometheus.enabled" default:"false"`
	Port    int    `mapstructure:"port" cmdx:"metric.prometheus.port" default:"9090"`
	Path    string `mapstructure:"path" cmdx:"metric.prometheus.path" default:"/metrics"`
}

type metric struct {
	RuntimeStatsRecordIntervalMS int64               `mapstructure:"runtime_stats_record_interval_ms" cmdx:"metric.runtime.stats.record.interval.ms" default:"10000"`
	StatsD                       metricStatsdCfg     `mapstructure:"statsd"`
	Prometheus                   metricPrometheusCfg `mapstructure:"prometheus"`
}
