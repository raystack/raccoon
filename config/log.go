package config

import (
	"github.com/odpf/raccoon/config/util"

	"github.com/spf13/viper"
)

var Log log

type log struct {
	Level string
}

func logConfigLoader() {
	viper.SetDefault("LOG_LEVEL", "info")
	Log = log{Level: util.MustGetString("LOG_LEVEL")}
}
