package config

import (
	"raccoon/config/util"

	"github.com/spf13/viper"
)

var Topic topic

type topic struct {
	Format string
}

func topicConfigLoader() {
	viper.SetDefault("TOPIC_FORMAT", "%s")
	Topic = topic{
		Format: util.MustGetString("TOPIC_FORMAT"),
	}
}
