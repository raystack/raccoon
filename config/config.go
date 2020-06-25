package config

import (
	"bytes"
	"github.com/spf13/viper"
)

var loaded bool

func Load() {
	if loaded {
		return
	}
	loaded = true
	viper.SetDefault("kafka_client_queue_buffering_max_messages", "100000")
	viper.SetDefault("kafka_flush_interval", "1000")
	viper.SetDefault("delivery_channel_size", "100")
	viper.SetConfigName("application")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	viper.MergeConfig(bytes.NewBuffer(dynamicKafkaConfigLoad()))
	viper.AutomaticEnv()

	ServerConfigLoader()
	LogLevel()
	NewKafkaConfig()
	WorkerConfigLoader()
}
