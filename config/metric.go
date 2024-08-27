package config

var Metric metric

type metricStatsdCfg struct {
	Enabled       bool   `mapstructure:"enabled" cmdx:"metric.statsd.enabled" default:"false"  desc:"Enable statsd metric exporter"`
	Address       string `mapstructure:"address" cmdx:"metric.statsd.address" default:":8125" desc:"Address to reports the service metrics"`
	FlushPeriodMS int64  `mapstructure:"flush_period_ms" cmdx:"metric.statsd.flush.period.ms" default:"10000" desc:"Interval for the service to push metrics"`
}

type metricPrometheusCfg struct {
	Enabled bool   `mapstructure:"enabled" cmdx:"metric.prometheus.enabled" default:"false" desc:"Enable prometheus http server to expose service metrics"`
	Port    int    `mapstructure:"port" cmdx:"metric.prometheus.port" default:"9090" desc:"Port to expose prometheus metrics on"`
	Path    string `mapstructure:"path" cmdx:"metric.prometheus.path" default:"/metrics" desc:"The path at which prometheus server should serve metrics"`
}

type metric struct {
	RuntimeStatsRecordIntervalMS int64               `mapstructure:"runtime_stats_record_interval_ms" cmdx:"metric.runtime.stats.record.interval.ms" default:"10000" desc:"Time interval between runtime metric collection"`
	StatsD                       metricStatsdCfg     `mapstructure:"statsd"`
	Prometheus                   metricPrometheusCfg `mapstructure:"prometheus"`
}
