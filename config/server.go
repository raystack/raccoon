package config

var Server = server{
	CORS: serverCors{
		// go-defaults doesn't support populating slice variables
		AllowedMethods: []string{"GET", "HEAD", "POST", "OPTIONS"},
	},
}

type server struct {
	CORS      serverCors `mapstructure:"cors"`
	GRPC      serverGRPC `mapstructure:"grpc"`
	Websocket serverWs   `mapstructure:"websocket"`
	Batch     batch      `mapstructure:"batch"`
}

type conn struct {
	IDHeader     string `mapstructure:"id_header" cmdx:"server.websocket.conn.id.header" desc:"Unique identifier for the server to maintain the connection"`
	GroupHeader  string `mapstructure:"group_header" cmdx:"server.websocket.conn.group.header" desc:"Additional identifier for the server to maintain the connection"`
	GroupDefault string `mapstructure:"group_default" cmdx:"server.websocket.conn.group.default" default:"--default--" desc:"Default connection group name"`
}

type batch struct {
	DedupEnabled bool `mapstructure:"dedup_in_connection_enabled" cmdx:"server.batch.dedup.in.connection.enabled" desc:"Whether to discard duplicate messages"`
}

type serverWs struct {
	Conn                conn   `mapstructure:"conn"`
	AppPort             string `mapstructure:"port" cmdx:"server.websocket.port" default:"8080" desc:"Port for the service to listen"`
	ServerMaxConn       int    `mapstructure:"max_conn" cmdx:"server.websocket.max.conn" default:"30000" desc:"Maximum connection that can be handled by the server instance"`
	ReadBufferSize      int    `mapstructure:"read_buffer_size" cmdx:"server.websocket.read.buffer.size" default:"10240" desc:"Input buffer size in bytes"`
	WriteBufferSize     int    `mapstructure:"write_buffer_size" cmdx:"server.websocket.write.buffer.size" default:"10240" desc:"Output buffer size in bytes"`
	PingIntervalMS      int64  `mapstructure:"ping_interval_ms" cmdx:"server.websocket.ping.interval.ms" default:"30000" desc:"Interval of each ping to client in milliseconds"`
	PongWaitIntervalMS  int64  `mapstructure:"pong_wait_interval_ms" cmdx:"server.websocket.pong.wait.interval.ms" default:"60000" desc:"Wait time for client to send Pong message in milliseconds"`
	WriteWaitIntervalMS int64  `mapstructure:"write_wait_interval_ms" cmdx:"server.websocket.write.wait.interval.ms" default:"5000" desc:"Timeout deadline set on the writes in milliseconds"`
	PingerSize          int    `mapstructure:"pinger_size" cmdx:"server.websocket.pinger.size" default:"1" desc:"Number of goroutine spawned to ping clients"`
	CheckOrigin         bool   `mapstructure:"check_origin" cmdx:"server.websocket.check.origin" default:"true" desc:"Toggle CORS check on WebSocket API"`
}

type serverGRPC struct {
	Port string `mapstructure:"port" cmdx:"server.grpc.port" default:"8081"`
}

type serverCors struct {
	Enabled          bool     `mapstructure:"enabled" cmdx:"server.cors.enabled" default:"false" desc:"Toggle CORS check on REST API"`
	AllowedOrigin    []string `mapstructure:"allowed_origin" cmdx:"server.cors.allowed.origin" desc:"Allowed origins for CORS. Use '*' to allow all"`
	AllowedMethods   []string `mapstructure:"allowed_methods" cmdx:"server.cors.allowed.methods" desc:"Allowed HTTP Methods for CORS"`
	AllowedHeaders   []string `mapstructure:"allowed_headers" cmdx:"server.cors.allowed.headers" desc:"Allowed HTTP Headers for CORS"`
	AllowCredentials bool     `mapstructure:"allow_credentials" cmdx:"server.cors.allow.credentials" default:"false" desc:"used to specify that the user agent may pass authentication details along with the request"`
	MaxAge           int      `mapstructure:"preflight_max_age_seconds" cmdx:"server.cors.preflight.max.age.seconds" desc:"Max Age of preflight responses"`
}
