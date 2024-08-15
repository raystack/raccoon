package config

var Server server

// type server struct {
// 	Publisher         string `mapstructure:"PUBLISHER_TYPE" cmdx:"publisher.type" default:"kafka"`
// 	PublisherKafka    publisherKafka
// 	PublisherPubSub   publisherPubSub
// 	PublisherKinesis  publisherKinesis
// 	Log               log
// 	Event             event
// 	EventDistribution eventDistribution
// 	MetricStatsd      metricStatsdCfg
// 	MetricPrometheus  metricPrometheusCfg
// 	MetricInfo        metricInfoCfg
// 	Websocket         serverWs
// 	CORS              serverCors
// 	GRPC              serverGRPC
// 	Worker            worker
// }

type server struct {
	CORS      serverCors `mapstructure:"cors"`
	GRPC      serverGRPC `mapstructure:"grpc"`
	Websocket serverWs   `mapstructure:"websocket"`
	Batch     batch      `mapstructure:"batch"`
}

// func (srv server) validate() error {
// 	trim := strings.TrimSpace
// 	if trim(srv.Websocket.ConnIDHeader) == "" {
// 		return errFieldRequired(srv.Websocket, "ConnIDHeader")
// 	}
// 	if srv.Publisher == "pubsub" {
// 		if trim(srv.PublisherPubSub.ProjectId) == "" {
// 			return errFieldRequired(srv.PublisherPubSub, "ProjectId")
// 		}
// 		if trim(srv.PublisherPubSub.CredentialsFile) == "" {
// 			return errFieldRequired(srv.PublisherPubSub, "CredentialsFile")
// 		}
// 	}

// 	if srv.Publisher == "kinesis" {

// 		hasAWSEnvCreds := trim(os.Getenv("AWS_ACCESS_KEY_ID")) != "" &&
// 			trim(os.Getenv("AWS_SECRET_ACCESS_KEY")) != ""

// 		if trim(srv.PublisherKinesis.CredentialsFile) == "" && !hasAWSEnvCreds {
// 			return errFieldRequired(srv.PublisherKinesis, "CredentialsFile")
// 		}
// 	}

// 	// there are no concrete fields that refer to this config
// 	kafkaServers := "PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS"
// 	if srv.Publisher == "kafka" && !viper.IsSet(kafkaServers) {
// 		flag := strings.ToLower(kafkaServers)
// 		flag = strings.ReplaceAll(flag, "_", ".")
// 		return errRequired(kafkaServers, flag)
// 	}

// 	validPublishers := []string{
// 		"kafka",
// 		"kinesis",
// 		"pubsub",
// 	}
// 	if !slices.Contains(validPublishers, srv.Publisher) {
// 		return fmt.Errorf("unknown publisher: %s", srv.Publisher)
// 	}

// 	return nil
// }

// prepare applies defaults and fallback values to Server.
// prepare must be called after loading configs using viper
// func prepare(srv *server) {

// 	// parse kafka dynamic config
// 	viper.MergeConfig(bytes.NewBuffer(dynamicKafkaClientConfigLoad()))

// 	// add fallback for pubsub credentials
// 	if srv.Publisher == "pubsub" {
// 		defaultCredentials := strings.TrimSpace(
// 			os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
// 		)
// 		if strings.TrimSpace(srv.PublisherPubSub.CredentialsFile) == "" && defaultCredentials != "" {
// 			srv.PublisherPubSub.CredentialsFile = defaultCredentials
// 		}
// 	}

// 	// add default CORS headers
// 	corsHeaders := []string{"Content-Type"}
// 	provisionalHeaders := []string{
// 		"SERVER_WEBSOCKET_CONN_GROUP_HEADER",
// 		"SERVER_WEBSOCKET_CONN_ID_HEADER",
// 	}

// 	for _, header := range provisionalHeaders {
// 		value := strings.TrimSpace(os.Getenv(header))
// 		if value != "" {
// 			corsHeaders = append(corsHeaders, value)
// 		}
// 	}

// 	for _, header := range corsHeaders {
// 		if !slices.Contains(srv.CORS.AllowedHeaders, header) {
// 			srv.CORS.AllowedHeaders = append(srv.CORS.AllowedHeaders, header)
// 		}
// 	}

// }

type conn struct {
	IDHeader     string `mapstructure:"id_header" cmdx:"server.websocket.conn.id.header"`
	GroupHeader  string `mapstructure:"group_header" cmdx:"server.websocket.conn.group.header"`
	GroupDefault string `mapstructure:"group_default" cmdx:"server.websocket.conn.group.default"`
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
	AllowedMethods   []string `mapstructure:"allowed_methods" cmdx:"server.cors.allowed.methods" default:"GET,HEAD,POST,OPTIONS"`
	AllowedHeaders   []string `mapstructure:"allowed_headers" cmdx:"server.cors.allowed.headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials" cmdx:"server.cors.allow.credentials" default:"false"`
	MaxAge           int      `mapstructure:"preflight_max_age_seconds" cmdx:"server.cors.preflight.max.age.seconds"`
}
