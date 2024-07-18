package config

import (
	"time"

	"github.com/raystack/raccoon/config/util"
	"github.com/spf13/viper"
)

var Server server
var ServerWs serverWs
var ServerGRPC serverGRPC
var ServerCors serverCors

type server struct {
	DedupEnabled bool
}

type serverWs struct {
	AppPort           string
	ServerMaxConn     int
	ReadBufferSize    int
	WriteBufferSize   int
	PingInterval      time.Duration
	PongWaitInterval  time.Duration
	WriteWaitInterval time.Duration
	PingerSize        int
	ConnIDHeader      string
	ConnGroupHeader   string
	ConnGroupDefault  string
	CheckOrigin       bool
}

type serverGRPC struct {
	Port string
}

type serverCors struct {
	Enabled          bool
	AllowedOrigin    []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func serverConfigLoader() {
	viper.SetDefault("SERVER_BATCH_DEDUP_IN_CONNECTION_ENABLED", "false")
	Server = server{
		DedupEnabled: util.MustGetBool("SERVER_BATCH_DEDUP_IN_CONNECTION_ENABLED"),
	}
}

func serverCorsConfigLoader() {
	viper.SetDefault("SERVER_CORS_ENABLED", false)
	viper.SetDefault("SERVER_CORS_ALLOWED_ORIGIN", "")
	viper.SetDefault("SERVER_CORS_ALLOWED_METHODS", []string{"GET", "HEAD", "POST", "OPTIONS"})
	viper.SetDefault("SERVER_CORS_ALLOWED_HEADERS", "")
	viper.SetDefault("SERVER_CORS_ALLOW_CREDENTIALS", false)
	viper.SetDefault("SERVER_CORS_PREFLIGHT_MAX_AGE_SECONDS", 0)
	ServerCors = serverCors{
		Enabled:          util.MustGetBool("SERVER_CORS_ENABLED"),
		AllowedOrigin:    viper.GetStringSlice("SERVER_CORS_ALLOWED_ORIGIN"),
		AllowedMethods:   viper.GetStringSlice("SERVER_CORS_ALLOWED_METHODS"),
		AllowCredentials: util.MustGetBool("SERVER_CORS_ALLOW_CREDENTIALS"),
		AllowedHeaders:   getAllowedHeaders(),
		MaxAge:           util.MustGetInt("SERVER_CORS_PREFLIGHT_MAX_AGE_SECONDS"),
	}
}

func getAllowedHeaders() []string {
	allowedHeaders := []string{"Content-Type"}
	allowedHeaders = setAllowedHeaders(allowedHeaders, "SERVER_WEBSOCKET_CONN_GROUP_HEADER")
	allowedHeaders = setAllowedHeaders(allowedHeaders, "SERVER_WEBSOCKET_CONN_ID_HEADER")
	inputHeaders := viper.GetStringSlice("SERVER_CORS_ALLOWED_HEADERS")
	for _, input := range inputHeaders {
		if input != "" && !util.Contains(input, allowedHeaders) {
			allowedHeaders = append(allowedHeaders, input)
		}
	}
	return allowedHeaders
}

func setAllowedHeaders(allowedHeaders []string, envKey string) []string {
	if header := viper.GetString(envKey); header != "" {
		allowedHeaders = append(allowedHeaders, header)
	}
	return allowedHeaders
}

func serverWsConfigLoader() {
	viper.SetDefault("SERVER_WEBSOCKET_CHECK_ORIGIN", true)
	viper.SetDefault("SERVER_WEBSOCKET_PORT", "8080")
	viper.SetDefault("SERVER_WEBSOCKET_MAX_CONN", 30000)
	viper.SetDefault("SERVER_WEBSOCKET_READ_BUFFER_SIZE", 10240)
	viper.SetDefault("SERVER_WEBSOCKET_WRITE_BUFFER_SIZE", 10240)
	viper.SetDefault("SERVER_WEBSOCKET_PING_INTERVAL_MS", "30000")
	viper.SetDefault("SERVER_WEBSOCKET_PONG_WAIT_INTERVAL_MS", "60000") //should be more than the ping period
	viper.SetDefault("SERVER_WEBSOCKET_WRITE_WAIT_INTERVAL_MS", "5000")
	viper.SetDefault("SERVER_WEBSOCKET_PINGER_SIZE", 1)
	viper.SetDefault("SERVER_WEBSOCKET_CONN_GROUP_HEADER", "")
	viper.SetDefault("SERVER_WEBSOCKET_CONN_GROUP_DEFAULT", "--default--")

	ServerWs = serverWs{
		AppPort:           util.MustGetString("SERVER_WEBSOCKET_PORT"),
		ServerMaxConn:     util.MustGetInt("SERVER_WEBSOCKET_MAX_CONN"),
		ReadBufferSize:    util.MustGetInt("SERVER_WEBSOCKET_READ_BUFFER_SIZE"),
		WriteBufferSize:   util.MustGetInt("SERVER_WEBSOCKET_WRITE_BUFFER_SIZE"),
		PingInterval:      util.MustGetDuration("SERVER_WEBSOCKET_PING_INTERVAL_MS", time.Millisecond),
		PongWaitInterval:  util.MustGetDuration("SERVER_WEBSOCKET_PONG_WAIT_INTERVAL_MS", time.Millisecond),
		WriteWaitInterval: util.MustGetDuration("SERVER_WEBSOCKET_WRITE_WAIT_INTERVAL_MS", time.Millisecond),
		PingerSize:        util.MustGetInt("SERVER_WEBSOCKET_PINGER_SIZE"),
		ConnIDHeader:      util.MustGetString("SERVER_WEBSOCKET_CONN_ID_HEADER"),
		ConnGroupHeader:   util.MustGetString("SERVER_WEBSOCKET_CONN_GROUP_HEADER"),
		ConnGroupDefault:  util.MustGetString("SERVER_WEBSOCKET_CONN_GROUP_DEFAULT"),
		CheckOrigin:       util.MustGetBool("SERVER_WEBSOCKET_CHECK_ORIGIN"),
	}
}

func serverGRPCConfigLoader() {
	viper.SetDefault("SERVER_GRPC_PORT", "8081")
	ServerGRPC = serverGRPC{
		Port: util.MustGetString("SERVER_GRPC_PORT"),
	}
}
