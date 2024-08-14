package config

import (
	"github.com/raystack/raccoon/config/util"

	"github.com/spf13/viper"
)

var Log log

type log struct {
	Level string `mapstructure:"LOG_LEVEL" cmdx:"log.level" default:"info" `
}

func logConfigLoader() {
	viper.SetDefault("LOG_LEVEL", "info")
	Log = log{Level: util.MustGetString("LOG_LEVEL")}
}
