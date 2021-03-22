package config

import (
	"raccoon/config/util"

	"github.com/spf13/viper"
)

type TopicConfig struct {
	Format string
}

func (tc TopicConfig) GetFormat() string {
	return tc.Format
}

func NewTopicConfig() TopicConfig {
	viper.SetDefault("TOPIC_FORMAT", "%s")
	tc := TopicConfig{
		Format: util.MustGetString("TOPIC_FORMAT"),
	}

	return tc
}
