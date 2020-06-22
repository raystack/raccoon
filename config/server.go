package config

import "github.com/spf13/viper"

type ServerConfig struct {
	AppPort         string
	ServerMaxConn   int
	ReadBufferSize  int
	WriteBufferSize int
	CheckOrigin     bool
}

func ServerConfigLoader() ServerConfig {
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("SERVER_MAX_CONN", 30000)
	viper.SetDefault("READ_BUFFER_SIZE", 10240)
	viper.SetDefault("WRITE_BUFFER_SIZE", 10240)
	viper.SetDefault("CHECK_ORIGIN", true)
	return ServerConfig{
		AppPort:         mustGetString("APP_PORT"),
		ServerMaxConn:   mustGetInt("SERVER_MAX_CONN"),
		ReadBufferSize:  mustGetInt("READ_BUFFER_SIZE"),
		WriteBufferSize: mustGetInt("WRITE_BUFFER_SIZE"),
		CheckOrigin:     mustGetBool("CHECK_ORIGIN"),
	}
}
