package config

import (
	"github.com/odpf/raccoon/config/util"
	"github.com/spf13/viper"
)

var Event event

type event struct {
	Ack int
}

func eventConfigLoader() {
	viper.SetDefault("EVENT_ACK", 0)
	Event = event{
		Ack: util.MustGetInt("EVENT_ACK"),
	}
}
