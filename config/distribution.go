package config

import (
	"github.com/raystack/raccoon/config/util"

	"github.com/spf13/viper"
)

var EventDistribution eventDistribution

type eventDistribution struct {
	PublisherPattern string `mapstructure:"EVENT_DISTRIBUTION_PUBLISHER_PATTERN" cmdx:"event.distribution.publisher.pattern"  default:"clickstream-%s-log"`
}

func eventDistributionConfigLoader() {
	viper.SetDefault("EVENT_DISTRIBUTION_PUBLISHER_PATTERN", "clickstream-%s-log")
	EventDistribution = eventDistribution{
		PublisherPattern: util.MustGetString("EVENT_DISTRIBUTION_PUBLISHER_PATTERN"),
	}
}
