package config

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	viper.Reset()
	viper.AutomaticEnv()
	os.Exit(m.Run())
}

func TestLogLevel(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	logConfigLoader()
	assert.Equal(t, "debug", Log.Level)
}

func TestServerConfig(t *testing.T) {
	// default value test
	serverConfigLoader()
	assert.Equal(t, false, Server.DedupEnabled)

	// override value test
	os.Setenv("SERVER_BATCH_DEDUP_IN_CONNECTION_ENABLED", "true")
	serverConfigLoader()
	assert.Equal(t, true, Server.DedupEnabled)
}

func TestServerWsConfig(t *testing.T) {
	os.Setenv("SERVER_WEBSOCKET_PORT", "8080")
	os.Setenv("SERVER_WEBSOCKET_PING_INTERVAL_MS", "1")
	os.Setenv("SERVER_WEBSOCKET_PONG_WAIT_INTERVAL_MS", "1")
	os.Setenv("SERVER_WEBSOCKET_SERVER_SHUTDOWN_GRACE_PERIOD_MS", "3")
	os.Setenv("SERVER_WEBSOCKET_CONN_ID_HEADER", "X-User-ID")
	serverWsConfigLoader()
	assert.Equal(t, "8080", ServerWs.AppPort)
	assert.Equal(t, time.Duration(1)*time.Millisecond, ServerWs.PingInterval)
	assert.Equal(t, time.Duration(1)*time.Millisecond, ServerWs.PongWaitInterval)
}

func TestGRPCServerConfig(t *testing.T) {
	os.Setenv("SERVER_GRPC_PORT", "8081")
	serverGRPCConfigLoader()
	assert.Equal(t, "8081", ServerGRPC.Port)
}

func TestDynamicConfigLoad(t *testing.T) {
	os.Setenv("PUBLISHER_KAFKA_CLIENT_RANDOM", "anything")
	os.Setenv("PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS", "localhost:9092")
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(dynamicKafkaClientConfigLoad()))
	assert.Equal(t, "anything", viper.GetString("PUBLISHER_KAFKA_CLIENT_RANDOM"))
	assert.Equal(t, "localhost:9092", viper.GetString("PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS"))
}

func TestKafkaConfig_ToKafkaConfigMap(t *testing.T) {
	os.Setenv("PUBLISHER_KAFKA_FLUSH_INTERVAL_MS", "1000")
	os.Setenv("PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS", "kafka:9092")
	os.Setenv("PUBLISHER_KAFKA_CLIENT_ACKS", "1")
	os.Setenv("PUBLISHER_KAFKA_CLIENT_QUEUE_BUFFERING_MAX_MESSAGES", "10000")
	os.Setenv("SOMETHING_PUBLISHER_KAFKA_CLIENT_SOMETHING", "anything")
	publisherKafkaConfigLoader()
	kafkaConfig := PublisherKafka.ToKafkaConfigMap()
	bootstrapServer, _ := kafkaConfig.Get("bootstrap.servers", "")
	topic, _ := kafkaConfig.Get("topic", "")
	something, _ := kafkaConfig.Get("client.something", "")
	assert.Equal(t, "kafka:9092", bootstrapServer)
	assert.Equal(t, "", topic)
	assert.NotEqual(t, something, "anything")
	assert.Equal(t, 4, len(*kafkaConfig))
}

func TestWorkerConfig(t *testing.T) {
	os.Setenv("WORKER_POOL_SIZE", "2")
	os.Setenv("WORKER_BUFFER_CHANNEL_SIZE", "5")
	os.Setenv("WORKER_KAFKA_DELIVERY_CHANNEL_SIZE", "10")
	os.Setenv("WORKER_BUFFER_FLUSH_TIMEOUT_MS", "100000")
	workerConfigLoader()
	assert.Equal(t, time.Duration(100)*time.Second, Worker.WorkerFlushTimeout)
	assert.Equal(t, 10, Worker.DeliveryChannelSize)
	assert.Equal(t, 5, Worker.ChannelSize)
	assert.Equal(t, 2, Worker.WorkersPoolSize)
}
