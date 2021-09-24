package config

import (
	"raccoon/config/util"
	"time"

	"github.com/spf13/viper"
)

var ServerWs serverWs

type serverWs struct {
	AppPort           string
	ServerMaxConn     int
	ReadBufferSize    int
	WriteBufferSize   int
	CheckOrigin       bool
	PingInterval      time.Duration
	PongWaitInterval  time.Duration
	WriteWaitInterval time.Duration
	PingerSize        int
	ConnIDHeader      string
	ConnTypeHeader    string
}

func serverWsConfigLoader() {
	viper.SetDefault("SERVER_WEBSOCKET_PORT", "8080")
	viper.SetDefault("SERVER_WEBSOCKET_MAX_CONN", 30000)
	viper.SetDefault("SERVER_WEBSOCKET_READ_BUFFER_SIZE", 10240)
	viper.SetDefault("SERVER_WEBSOCKET_WRITE_BUFFER_SIZE", 10240)
	viper.SetDefault("SERVER_WEBSOCKET_CHECK_ORIGIN", true)
	viper.SetDefault("SERVER_WEBSOCKET_PING_INTERVAL_MS", "30000")
	viper.SetDefault("SERVER_WEBSOCKET_PONG_WAIT_INTERVAL_MS", "60000") //should be more than the ping period
	viper.SetDefault("SERVER_WEBSOCKET_WRITE_WAIT_INTERVAL_MS", "5000")
	viper.SetDefault("SERVER_WEBSOCKET_PINGER_SIZE", 1)
	viper.SetDefault("SERVER_WEBSOCKET_CONN_TYPE_HEADER", "")

	ServerWs = serverWs{
		AppPort:           util.MustGetString("SERVER_WEBSOCKET_PORT"),
		ServerMaxConn:     util.MustGetInt("SERVER_WEBSOCKET_MAX_CONN"),
		ReadBufferSize:    util.MustGetInt("SERVER_WEBSOCKET_READ_BUFFER_SIZE"),
		WriteBufferSize:   util.MustGetInt("SERVER_WEBSOCKET_WRITE_BUFFER_SIZE"),
		CheckOrigin:       util.MustGetBool("SERVER_WEBSOCKET_CHECK_ORIGIN"),
		PingInterval:      util.MustGetDuration("SERVER_WEBSOCKET_PING_INTERVAL_MS", time.Millisecond),
		PongWaitInterval:  util.MustGetDuration("SERVER_WEBSOCKET_PONG_WAIT_INTERVAL_MS", time.Millisecond),
		WriteWaitInterval: util.MustGetDuration("SERVER_WEBSOCKET_WRITE_WAIT_INTERVAL_MS", time.Microsecond),
		PingerSize:        util.MustGetInt("SERVER_WEBSOCKET_PINGER_SIZE"),
		ConnIDHeader:      util.MustGetString("SERVER_WEBSOCKET_CONN_ID_HEADER"),
		ConnTypeHeader:    util.MustGetString("SERVER_WEBSOCKET_CONN_TYPE_HEADER"),
	}
}
