package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Statsd struct {
	Address       string
	FlushPeriodMs int
}

func (s Statsd) FlushPeriod() time.Duration {
	d, err := time.ParseDuration(fmt.Sprintf("%dms", s.FlushPeriodMs))
	if err != nil {
		panic(fmt.Sprintf("FlushPeriod cannot be parsed: %v", err))
	}
	return d
}

func StatsdConfigLoader() Statsd {
	viper.SetDefault("STATSD_ADDRESS", ":8125")
	viper.SetDefault("STATSD_FLUSH_PERIOD_MS", 10000)
	return Statsd{
		Address:       mustGetString("STATSD_ADDRESS"),
		FlushPeriodMs: mustGetInt("STATSD_FLUSH_PERIOD_MS"),
	}
}
