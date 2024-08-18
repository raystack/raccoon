package config

var Log log

type log struct {
	Level string `mapstructure:"level" cmdx:"log.level" default:"info" desc:"Levels available are [debug info warn error fatal panic]"`
}
