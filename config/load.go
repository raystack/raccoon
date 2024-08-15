package config

import (
	"errors"

	defaults "github.com/mcuadros/go-defaults"
	"github.com/raystack/salt/config"
	"github.com/spf13/viper"
)

// Load configs from env or yaml and set it to respective keys
func Load() error {
	conf := cfg{&Server, &Publisher, &Worker, &Event, &Metric, &Log}
	loader := config.NewLoader(
		config.WithViper(viper.GetViper()),
		config.WithName(".env"),
		config.WithPath("./"),
		config.WithPath("../"),
		config.WithPath("../../"),
		config.WithType("env"),
	)
	err := loader.Load(&conf)
	if err != nil && !errors.As(err, &config.ConfigFileNotFoundError{}) {
		return err
	}

	prepare()
	return validate()
}

func init() {
	defaults.SetDefaults(&Server)
}
