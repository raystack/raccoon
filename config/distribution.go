package config

import (
	"raccoon/config/util"

	"github.com/spf13/viper"
)

var DISTRIBUTION distribution

type distribution struct {
	PublisherPattern string
}

func distributionConfigLoader() {
	viper.SetDefault("DISTRIBUTION-EVENT_DISTRIBUTION_PUBLISHER_PATTERN", "clickstream-%s-log")
	DISTRIBUTION = distribution{
		PublisherPattern: util.MustGetString("DISTRIBUTION-EVENT_DISTRIBUTION_PUBLISHER_PATTERN"),
	}
}
