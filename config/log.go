package config

import "github.com/spf13/viper"

func LogLevel() string {
	viper.SetDefault("LOG_LEVEL", "info")
	return mustGetString("LOG_LEVEL")
}
