package config

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/spf13/viper"
)

var Server server

type server struct {
	Publisher         string `mapstructure:"PUBLISHER_TYPE" cmdx:"publisher.type" default:"kafka"`
	PublisherKafka    publisherKafka
	PublisherPubSub   publisherPubSub
	PublisherKinesis  publisherKinesis
	Log               log
	Event             event
	EventDistribution eventDistribution
	MetricStatsd      metricStatsdCfg
	MetricPrometheus  metricPrometheusCfg
	MetricInfo        metricInfoCfg
	Websocket         serverWs
	CORS              serverCors
	GRPC              serverGRPC
	Worker            worker
}

func (srv server) validate() error {
	trim := strings.TrimSpace
	if trim(srv.Websocket.ConnIDHeader) == "" {
		return errFieldRequired(srv.Websocket, "ConnIDHeader")
	}
	if srv.Publisher == "pubsub" {
		if trim(srv.PublisherPubSub.ProjectId) == "" {
			return errFieldRequired(srv.PublisherPubSub, "ProjectId")
		}
		if trim(srv.PublisherPubSub.CredentialsFile) == "" {
			return errFieldRequired(srv.PublisherPubSub, "CredentialsFile")
		}
	}

	if srv.Publisher == "kinesis" {

		hasAWSEnvCreds := trim(os.Getenv("AWS_ACCESS_KEY_ID")) != "" &&
			trim(os.Getenv("AWS_SECRET_ACCESS_KEY")) != ""

		if trim(srv.PublisherKinesis.CredentialsFile) == "" && !hasAWSEnvCreds {
			return errFieldRequired(srv.PublisherKinesis, "CredentialsFile")
		}
	}

	// there are no concrete fields that refer to this config
	kafkaServers := "PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS"
	if srv.Publisher == "kafka" && !viper.IsSet(kafkaServers) {
		flag := strings.ToLower(kafkaServers)
		flag = strings.ReplaceAll(flag, "_", ".")
		return errRequired(kafkaServers, flag)
	}

	validPublishers := []string{
		"kafka",
		"kinesis",
		"pubsub",
	}
	if !slices.Contains(validPublishers, srv.Publisher) {
		return fmt.Errorf("unknown publisher: %s", srv.Publisher)
	}

	return nil
}

// prepare applies defaults and fallback values to Server.
// prepare must be called after loading configs using viper
func prepare(srv *server) {

	// parse kafka dynamic config
	viper.MergeConfig(bytes.NewBuffer(dynamicKafkaClientConfigLoad()))

	// add fallback for pubsub credentials
	if srv.Publisher == "pubsub" {
		defaultCredentials := strings.TrimSpace(
			os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		)
		if strings.TrimSpace(srv.PublisherPubSub.CredentialsFile) == "" && defaultCredentials != "" {
			srv.PublisherPubSub.CredentialsFile = defaultCredentials
		}
	}

	// add default CORS headers
	corsHeaders := []string{"Content-Type"}
	provisionalHeaders := []string{
		"SERVER_WEBSOCKET_CONN_GROUP_HEADER",
		"SERVER_WEBSOCKET_CONN_ID_HEADER",
	}

	for _, header := range provisionalHeaders {
		value := strings.TrimSpace(os.Getenv(header))
		if value != "" {
			corsHeaders = append(corsHeaders, value)
		}
	}

	for _, header := range corsHeaders {
		if !slices.Contains(srv.CORS.AllowedHeaders, header) {
			srv.CORS.AllowedHeaders = append(srv.CORS.AllowedHeaders, header)
		}
	}

}

type serverWs struct {
	AppPort             string `mapstructure:"SERVER_WEBSOCKET_PORT" cmdx:"server.websocket.port" default:"8080"`
	ServerMaxConn       int    `mapstructure:"SERVER_WEBSOCKET_MAX_CONN" cmdx:"server.websocket.max.conn" default:"30000"`
	ReadBufferSize      int    `mapstructure:"SERVER_WEBSOCKET_READ_BUFFER_SIZE" cmdx:"server.websocket.read.buffer.size" default:"10240"`
	WriteBufferSize     int    `mapstructure:"SERVER_WEBSOCKET_WRITE_BUFFER_SIZE" cmdx:"server.websocket.write.buffer.size" default:"10240"`
	PingIntervalMS      int64  `mapstructure:"SERVER_WEBSOCKET_PING_INTERVAL_MS" cmdx:"server.websocket.ping.interval.ms" default:"30000"`
	PongWaitIntervalMS  int64  `mapstructure:"SERVER_WEBSOCKET_PONG_WAIT_INTERVAL_MS" cmdx:"server.websocket.pong.wait.interval.ms" default:"60000"`
	WriteWaitIntervalMS int64  `mapstructure:"SERVER_WEBSOCKET_WRITE_WAIT_INTERVAL_MS" cmdx:"server.websocket.write.wait.interval.ms" default:"5000"`
	PingerSize          int    `mapstructure:"SERVER_WEBSOCKET_PINGER_SIZE" cmdx:"server.websocket.pinger.size" default:"1"`
	ConnIDHeader        string `mapstructure:"SERVER_WEBSOCKET_CONN_ID_HEADER" cmdx:"server.websocket.conn.id.header"`
	ConnGroupHeader     string `mapstructure:"SERVER_WEBSOCKET_CONN_GROUP_HEADER" cmdx:"server.websocket.conn.group.header"`
	ConnGroupDefault    string `mapstructure:"SERVER_WEBSOCKET_CONN_GROUP_DEFAULT" cmdx:"server.websocket.conn.group.default" default:"--default--"`
	CheckOrigin         bool   `mapstructure:"SERVER_WEBSOCKET_CHECK_ORIGIN" cmdx:"server.websocket.check.origin" default:"true"`
	DedupEnabled        bool   `mapstructure:"SERVER_BATCH_DEDUP_IN_CONNECTION_ENABLED" cmdx:"server.batch.dedup.in.connection.enabled" default:"false"`
}

type serverGRPC struct {
	Port string `mapstructure:"SERVER_GRPC_PORT" cmdx:"server.grpc.port" default:"8081"`
}

type serverCors struct {
	Enabled          bool     `mapstructure:"SERVER_CORS_ENABLED" cmdx:"server.cors.enabled" default:"false"`
	AllowedOrigin    []string `mapstructure:"SERVER_CORS_ALLOWED_ORIGIN" cmdx:"server.cors.allowed.origin"`
	AllowedMethods   []string `mapstructure:"SERVER_CORS_ALLOWED_METHODS" cmdx:"server.cors.allowed.methods" default:"GET,HEAD,POST,OPTIONS"`
	AllowedHeaders   []string `mapstructure:"SERVER_CORS_ALLOWED_HEADERS" cmdx:"server.cors.allowed.headers"`
	AllowCredentials bool     `mapstructure:"SERVER_CORS_ALLOW_CREDENTIALS" cmdx:"server.cors.allow.credentials" default:"false"`
	MaxAge           int      `mapstructure:"SERVER_CORS_PREFLIGHT_MAX_AGE_SECONDS" cmdx:"server.cors.preflight.max.age.seconds"`
}
