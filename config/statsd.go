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

func statsdConfigLoader() {
	viper.SetDefault("STATSD_ADDRESS", ":8125")
	viper.SetDefault("STATSD_FLUSH_PERIOD_MS", 10000)
	Statsd = statsd{
		Address:       util.MustGetString("STATSD_ADDRESS"),
		FlushPeriodMs: util.MustGetDuration("STATSD_FLUSH_PERIOD_MS", time.Millisecond),
	}
}
