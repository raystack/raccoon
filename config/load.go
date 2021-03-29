package config

import (
	"github.com/spf13/viper"
)

var loaded bool

// Load configs from env or yaml and set it to respective keys
func Load() {
	if loaded {
		return
	}
	loaded = true
	viper.AutomaticEnv()
	viper.SetConfigName("application")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	logConfigLoader()
	publisherKafkaConfigLoader()
	serverWsConfigLoader()
	workerConfigLoader()
	metricStatsdConfigLoader()
	eventDistributionConfigLoader()
}