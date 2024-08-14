package config

import (
	"bytes"

	defaults "github.com/mcuadros/go-defaults"
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

	viper.MergeConfig(bytes.NewBuffer(dynamicKafkaClientConfigLoad()))

}

func init() {
	defaults.SetDefaults(&Server)
}
