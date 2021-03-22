package config

import (
	"raccoon/config/util"

	"github.com/spf13/viper"
)

func LogLevel() string {
	viper.SetDefault("LOG_LEVEL", "info")
	return util.MustGetString("LOG_LEVEL")
}
