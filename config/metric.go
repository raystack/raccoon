package config

import (
	"time"

	"github.com/raystack/raccoon/config/util"

	"github.com/spf13/viper"
)

var MetricStatsd metricStatsdCfg

type metricStatsdCfg struct {
	Address       string
	FlushPeriodMs time.Duration
}

func metricStatsdConfigLoader() {
	viper.SetDefault("METRIC_STATSD_ADDRESS", ":8125")
	viper.SetDefault("METRIC_STATSD_FLUSH_PERIOD_MS", 10000)
	MetricStatsd = metricStatsdCfg{
		Address:       util.MustGetString("METRIC_STATSD_ADDRESS"),
		FlushPeriodMs: util.MustGetDuration("METRIC_STATSD_FLUSH_PERIOD_MS", time.Millisecond),
	}
}
