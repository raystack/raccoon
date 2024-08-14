package config

var Log log

type log struct {
	Level string `mapstructure:"LOG_LEVEL" cmdx:"log.level" default:"info" `
}
