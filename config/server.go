package config

import (
	"raccoon/config/util"
	"time"

	"github.com/spf13/viper"
)

var Websocket websocket

type websocket struct {
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
	viper.SetDefault("SERVER-WEBSOCKET-PORT", "8080")
	viper.SetDefault("SERVER-WEBSOCKET-MAX_CONN", 30000)
	viper.SetDefault("SERVER-WEBSOCKET-READ_BUFFER_SIZE", 10240)
	viper.SetDefault("SERVER-WEBSOCKET-WRITE_BUFFER_SIZE", 10240)
	viper.SetDefault("SERVER-WEBSOCKET-CHECK_ORIGIN", true)
	viper.SetDefault("SERVER-WEBSOCKET-PING_INTERVAL", "30")
	viper.SetDefault("SERVER-WEBSOCKET-PONG_WAIT_INTERVAL", "60") //should be more than the ping period
	viper.SetDefault("SERVER-WEBSOCKET-WRITE_WAIT_INTERVAL", "5")
	viper.SetDefault("SERVER-WEBSOCKET-SERVER_SHUTDOWN_GRACE_PERIOD", "3")
	viper.SetDefault("SERVER-WEBSOCKET-PINGER_SIZE", 1)

	Websocket = websocket{
		AppPort:                   util.MustGetString("SERVER-WEBSOCKET-PORT"),
		ServerMaxConn:             util.MustGetInt("SERVER-WEBSOCKET-MAX_CONN"),
		ReadBufferSize:            util.MustGetInt("SERVER-WEBSOCKET-READ_BUFFER_SIZE"),
		WriteBufferSize:           util.MustGetInt("SERVER-WEBSOCKET-WRITE_BUFFER_SIZE"),
		CheckOrigin:               util.MustGetBool("SERVER-WEBSOCKET-CHECK_ORIGIN"),
		PingInterval:              util.MustGetDuration("SERVER-WEBSOCKET-PING_INTERVAL", time.Second),
		PongWaitInterval:          util.MustGetDuration("SERVER-WEBSOCKET-PONG_WAIT_INTERVAL", time.Second),
		WriteWaitInterval:         util.MustGetDuration("SERVER-WEBSOCKET-WRITE_WAIT_INTERVAL", time.Second),
		ServerShutDownGracePeriod: util.MustGetDuration("SERVER-WEBSOCKET-SERVER_SHUTDOWN_GRACE_PERIOD", time.Second),
		PingerSize:                util.MustGetInt("SERVER-WEBSOCKET-PINGER_SIZE"),
		UserIDHeader:              util.MustGetString("SERVER-WEBSOKCET-USER_ID_HEADER"),
	}
}
