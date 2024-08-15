package config

var Event event

type AckType int

const (
	Asynchronous AckType = 0
	Synchronous  AckType = 1
)

type event struct {
	Ack                          AckType `mapstructure:"ack" cmdx:"event.ack" default:"0"`
	DistributionPublisherPattern string  `mapstructure:"distribution_publisher_pattern" cmdx:"event.distribution.publisher.pattern"  default:"clickstream-%s-log"`
}
