package config

import (
	"time"

	"github.com/spf13/viper"
)

var ServerConfig ServerCfg

type ServerCfg struct {
	AppPort           string
	ServerMaxConn     int
	ReadBufferSize    int
	WriteBufferSize   int
	CheckOrigin       bool
	PingInterval      time.Duration
	PongWaitInterval  time.Duration
	WriteWaitInterval time.Duration
}

func serverConfigLoader() {
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("SERVER_MAX_CONN", 30000)
	viper.SetDefault("READ_BUFFER_SIZE", 10240)
	viper.SetDefault("WRITE_BUFFER_SIZE", 10240)
	viper.SetDefault("CHECK_ORIGIN", true)
	viper.SetDefault("PING_INTERVAL", "30")
	viper.SetDefault("PONG_WAIT_INTERVAL", "60") //should be more than the ping period
	viper.SetDefault("WRITE_WAIT_INTERVAL", "5")

	ServerConfig = ServerCfg{
		AppPort:           mustGetString("APP_PORT"),
		ServerMaxConn:     mustGetInt("SERVER_MAX_CONN"),
		ReadBufferSize:    mustGetInt("READ_BUFFER_SIZE"),
		WriteBufferSize:   mustGetInt("WRITE_BUFFER_SIZE"),
		CheckOrigin:       mustGetBool("CHECK_ORIGIN"),
		PingInterval:      mustGetDurationInSeconds("PING_INTERVAL"),
		PongWaitInterval:  mustGetDurationInSeconds("PONG_WAIT_INTERVAL"),
		WriteWaitInterval: mustGetDurationInSeconds("WRITE_WAIT_INTERVAL"),
	}
}
