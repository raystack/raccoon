package config

var EventDistribution eventDistribution

type eventDistribution struct {
	PublisherPattern string `mapstructure:"EVENT_DISTRIBUTION_PUBLISHER_PATTERN" cmdx:"event.distribution.publisher.pattern"  default:"clickstream-%s-log"`
}
