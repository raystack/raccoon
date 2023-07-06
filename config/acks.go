package config

import (
	"github.com/raystack/raccoon/config/util"
	"github.com/spf13/viper"
)

var Event event

type AckType int

const (
	Asynchronous AckType = 0
	Synchronous  AckType = 1
)

type event struct {
	Ack AckType
}

func eventConfigLoader() {
	viper.SetDefault("EVENT_ACK", 0)
	Event = event{
		Ack: AckType(util.MustGetInt("EVENT_ACK")),
	}
}
