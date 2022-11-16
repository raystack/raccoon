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
	viper.SetConfigName(".env")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.SetConfigType("env")
	viper.ReadInConfig()

	logConfigLoader()
	publisherKafkaConfigLoader()
	serverConfigLoader()
	serverWsConfigLoader()
	serverGRPCConfigLoader()
	workerConfigLoader()
	metricStatsdConfigLoader()
	eventDistributionConfigLoader()
	eventConfigLoader()
}
