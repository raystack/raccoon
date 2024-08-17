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
	IDHeader     string `mapstructure:"id_header" cmdx:"server.websocket.conn.id.header"`
	GroupHeader  string `mapstructure:"group_header" cmdx:"server.websocket.conn.group.header"`
	GroupDefault string `mapstructure:"group_default" cmdx:"server.websocket.conn.group.default" default:"--default--"`
}

type batch struct {
	DedupEnabled bool `mapstructure:"dedup_in_connection_enabled" cmdx:"server.batch.dedup.in.connection.enabled"`
}

type serverWs struct {
	Conn                conn   `mapstructure:"conn"`
	AppPort             string `mapstructure:"port" cmdx:"server.websocket.port" default:"8080"`
	ServerMaxConn       int    `mapstructure:"max_conn" cmdx:"server.websocket.max.conn" default:"30000"`
	ReadBufferSize      int    `mapstructure:"read_buffer_size" cmdx:"server.websocket.read.buffer.size" default:"10240"`
	WriteBufferSize     int    `mapstructure:"write_buffer_size" cmdx:"server.websocket.write.buffer.size" default:"10240"`
	PingIntervalMS      int64  `mapstructure:"ping_interval_ms" cmdx:"server.websocket.ping.interval.ms" default:"30000"`
	PongWaitIntervalMS  int64  `mapstructure:"pong_wait_interval_ms" cmdx:"server.websocket.pong.wait.interval.ms" default:"60000"`
	WriteWaitIntervalMS int64  `mapstructure:"write_wait_interval_ms" cmdx:"server.websocket.write.wait.interval.ms" default:"5000"`
	PingerSize          int    `mapstructure:"pinger_size" cmdx:"server.websocket.pinger.size" default:"1"`
	CheckOrigin         bool   `mapstructure:"check_origin" cmdx:"server.websocket.check.origin" default:"true"`
}

type serverGRPC struct {
	Port string `mapstructure:"port" cmdx:"server.grpc.port" default:"8081"`
}

type serverCors struct {
	Enabled          bool     `mapstructure:"enabled" cmdx:"server.cors.enabled" default:"false"`
	AllowedOrigin    []string `mapstructure:"allowed_origin" cmdx:"server.cors.allowed.origin"`
	AllowedMethods   []string `mapstructure:"allowed_methods" cmdx:"server.cors.allowed.methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers" cmdx:"server.cors.allowed.headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials" cmdx:"server.cors.allow.credentials" default:"false"`
	MaxAge           int      `mapstructure:"preflight_max_age_seconds" cmdx:"server.cors.preflight.max.age.seconds"`
}
