package config

import (
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
		Format: mustGetString("TOPIC_FORMAT"),
	}

	return tc
}
