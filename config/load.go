package config

import (
	defaults "github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

// Load configs from env or yaml and set it to respective keys
func Load() error {
	viper.AutomaticEnv()
	viper.SetConfigName(".env")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.SetConfigType("env")
	viper.ReadInConfig()

	prepare(&Server)
	return Server.validate()
}

func init() {
	defaults.SetDefaults(&Server)
}
