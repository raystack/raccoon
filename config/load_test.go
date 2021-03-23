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
	os.Setenv("APP_PORT", "8080")
	os.Setenv("PING_INTERVAL", "1")
	os.Setenv("PONG_WAIT_INTERVAL", "1")
	os.Setenv("SERVER_SHUTDOWN_GRACE_PERIOD", "3")
	os.Setenv("USER_ID_HEADER", "x-user-id")
	serverConfigLoader()
	assert.Equal(t, "8080", Server.AppPort)
	assert.Equal(t, time.Duration(1)*time.Second, Server.PingInterval)
	assert.Equal(t, time.Duration(1)*time.Second, Server.PongWaitInterval)
	assert.Equal(t, time.Duration(3)*time.Second, Server.ServerShutDownGracePeriod)
}

func TestDynamicConfigLoad(t *testing.T) {
	os.Setenv("KAFKA_CLIENT_RANDOM", "anything")
	os.Setenv("KAFKA_CLIENT_BOOTSTRAP_SERVERS", "localhost:9092")
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(dynamicKafkaClientConfigLoad()))
	assert.Equal(t, "anything", viper.GetString("kafka_client_random"))
	assert.Equal(t, "localhost:9092", viper.GetString("kafka_client_bootstrap_servers"))
}

func TestKafkaConfig_ToKafkaConfigMap(t *testing.T) {
	os.Setenv("KAFKA_TOPIC", "test1")
	os.Setenv("KAFKA_FLUSH_INTERVAL", "1000")
	os.Setenv("KAFKA_CLIENT_BOOTSTRAP_SERVERS", "kafka:9092")
	os.Setenv("KAFKA_CLIENT_ACKS", "1")
	os.Setenv("KAFKA_CLIENT_QUEUE_BUFFERING_MAX_MESSAGES", "10000")
	os.Setenv("SOMETHING_KAFKA_CLIENT_SOMETHING", "anything")

	kafkaConfigLoader()
	kafkaConfig := Kafka.ToKafkaConfigMap()
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
	os.Setenv("BUFFER_CHANNEL_SIZE", "5")
	os.Setenv("DELIVERY_CHANNEL_SIZE", "10")
	os.Setenv("WORKER_FLUSH_TIMEOUT", "100")
	workerConfigLoader()
	assert.Equal(t, 100, Worker.WorkerFlushTimeout)
	assert.Equal(t, 10, Worker.DeliveryChannelSize)
	assert.Equal(t, 5, Worker.ChannelSize)
	assert.Equal(t, 2, Worker.WorkersPoolSize)
}
