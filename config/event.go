package config

var Event event

type AckType int

const (
	Asynchronous AckType = 0
	Synchronous  AckType = 1
)

type event struct {
	Ack                          AckType `mapstructure:"ack" cmdx:"event.ack" default:"0" desc:"Whether to send acknowledgements to clients or not. 1 to enable, 0 to disable."`
	DistributionPublisherPattern string  `mapstructure:"distribution_publisher_pattern" cmdx:"event.distribution.publisher.pattern"  default:"clickstream-%s-log" desc:"Topic template used for routing events"`
}
