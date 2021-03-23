package config

import (
	"raccoon/config/util"
	"time"

	"github.com/spf13/viper"
)

var Server server

type server struct {
	AppPort                   string
	ServerMaxConn             int
	ReadBufferSize            int
	WriteBufferSize           int
	CheckOrigin               bool
	PingInterval              time.Duration
	PongWaitInterval          time.Duration
	WriteWaitInterval         time.Duration
	ServerShutDownGracePeriod time.Duration
	PingerSize                int
	UserIDHeader              string
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
	viper.SetDefault("SERVER_SHUTDOWN_GRACE_PERIOD", "3")
	viper.SetDefault("PINGER_SIZE", 1)

	Server = server{
		AppPort:                   util.MustGetString("APP_PORT"),
		ServerMaxConn:             util.MustGetInt("SERVER_MAX_CONN"),
		ReadBufferSize:            util.MustGetInt("READ_BUFFER_SIZE"),
		WriteBufferSize:           util.MustGetInt("WRITE_BUFFER_SIZE"),
		CheckOrigin:               util.MustGetBool("CHECK_ORIGIN"),
		PingInterval:              util.MustGetDuration("PING_INTERVAL", time.Second),
		PongWaitInterval:          util.MustGetDuration("PONG_WAIT_INTERVAL", time.Second),
		WriteWaitInterval:         util.MustGetDuration("WRITE_WAIT_INTERVAL", time.Second),
		ServerShutDownGracePeriod: util.MustGetDuration("SERVER_SHUTDOWN_GRACE_PERIOD", time.Second),
		PingerSize:                util.MustGetInt("PINGER_SIZE"),
		UserIDHeader:              util.MustGetString("USER_ID_HEADER"),
	}
}
