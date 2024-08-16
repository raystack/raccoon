package config

import (
	"errors"

	defaults "github.com/mcuadros/go-defaults"
	"github.com/raystack/salt/config"
)

// Load configs from env or yaml and set it to respective keys
func Load(configFile string) error {
	loader := config.NewLoader(
		config.WithFile(configFile),
	)
	err := loader.Load(&cfg)
	if err != nil && !errors.As(err, &config.ConfigFileNotFoundError{}) {
		return err
	}

	prepare()
	return validate()
}

func init() {
	// go-defaults doesn't work properly with nested pointer values,
	// so we have to individually set defaults for each config class
	defaults.SetDefaults(&Server)
	defaults.SetDefaults(&Publisher)
	defaults.SetDefaults(&Worker)
	defaults.SetDefaults(&Event)
	defaults.SetDefaults(&Metric)
	defaults.SetDefaults(&Log)
}
