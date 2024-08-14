package config

var Event event

type AckType int

const (
	Asynchronous AckType = 0
	Synchronous  AckType = 1
)

type event struct {
	Ack AckType `mapstructure:"EVENT_ACK" cmdx:"event.ack" default:"0"`
}
