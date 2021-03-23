package config

import (
	"raccoon/config/util"
	"time"

	"github.com/spf13/viper"
)

var Statsd statsd

type statsd struct {
	Address       string
	FlushPeriodMs time.Duration
}

func metricConfigLoader() {
	viper.SetDefault("METRIC-STATSD-ADDRESS", ":8125")
	viper.SetDefault("METRIC-STATSD-FLUSH_PERIOD_MS", 10000)
	Statsd = statsd{
		Address:       util.MustGetString("METRIC-STATSD-ADDRESS"),
		FlushPeriodMs: util.MustGetDuration("METRIC-STATSD-FLUSH_PERIOD_MS", time.Millisecond),
	}
}
