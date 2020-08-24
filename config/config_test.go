package config

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLogLevel(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	viper.AutomaticEnv()
	assert.Equal(t, "debug", LogLevel())

	viper.Reset()
}

func TestServerConfig(t *testing.T) {
	os.Setenv("APP_PORT", "8080")
	os.Setenv("PING_INTERVAL", "1")
	os.Setenv("PONG_WAIT_INTERVAL", "1")
	os.Setenv("SERVER_SHUTDOWN_GRACE_PERIOD", "3")
	viper.AutomaticEnv()
	serverConfigLoader()
	assert.Equal(t, "8080", ServerConfig.AppPort)
	assert.Equal(t, time.Duration(1)*time.Second, ServerConfig.PingInterval)
	assert.Equal(t, time.Duration(1)*time.Second, ServerConfig.PongWaitInterval)
	assert.Equal(t, time.Duration(3)*time.Second, ServerConfig.ServerShutDownGracePeriod)

	viper.Reset()
}

func TestNewKafkaConfig(t *testing.T) {
	os.Setenv("TOPIC_FORMAT", "%s")
	os.Setenv("KAFKA_FLUSH_INTERVAL", "1000")

	expectedKafkaConfig := KafkaConfig{
		FlushInterval: 1000,
		TopicFormat: "%s",
	}

	viper.AutomaticEnv()
	kafkaConfig := NewKafkaConfig()
	assert.Equal(t, expectedKafkaConfig, kafkaConfig)

	viper.Reset()
}

func TestDynamicConfigLoad(t *testing.T) {
	os.Setenv("KAFKA_CLIENT_RANDOM", "anything")
	os.Setenv("KAFKA_CLIENT_BOOTSTRAP_SERVERS", "localhost:9092")
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(dynamicKafkaConfigLoad()))
	assert.Equal(t, "anything", viper.AllSettings()["kafka_client_random"])
	assert.Equal(t, "localhost:9092", viper.AllSettings()["kafka_client_bootstrap_servers"])

	viper.Reset()
}

func TestKafkaConfig_ToKafkaConfigMap(t *testing.T) {
	os.Setenv("KAFKA_TOPIC", "test1")
	os.Setenv("KAFKA_FLUSH_INTERVAL", "1000")
	os.Setenv("KAFKA_CLIENT_BOOTSTRAP_SERVERS", "kafka:6668")
	os.Setenv("KAFKA_CLIENT_ACKS", "1")
	os.Setenv("KAFKA_CLIENT_QUEUE_BUFFERING_MAX_MESSAGES", "10000")
	os.Setenv("SOMETHING_KAFKA_CLIENT_SOMETHING", "anything")

	viper.AutomaticEnv()
	viper.BindEnv("KAFKA_TOPIC")
	viper.BindEnv("KAFKA_FLUSH_INTERVAL")
	viper.BindEnv("KAFKA_CLIENT_BOOTSTRAP_SERVERS")
	viper.BindEnv("KAFKA_CLIENT_ACKS")
	viper.BindEnv("KAFKA_CLIENT_QUEUE_BUFFERING_MAX_MESSAGES")
	viper.BindEnv("SOMETHING_KAFKA_CLIENT_SOMETHING")

	kafkaConfig := NewKafkaConfig().ToKafkaConfigMap()
	bootstrapServer, _ := kafkaConfig.Get("bootstrap.servers", "")
	topic, _ := kafkaConfig.Get("topic", "")
	something, _ := kafkaConfig.Get("client.something", "")
	assert.Equal(t, "kafka:6668", bootstrapServer)
	assert.Equal(t, "", topic)
	assert.NotEqual(t, something, "anything")
	assert.Equal(t, 3, len(*kafkaConfig))

	viper.Reset()
}

func TestWorkerConfig(t *testing.T) {
	os.Setenv("WORKER_POOL_SIZE", "2")
	os.Setenv("BUFFER_CHANNEL_SIZE", "5")
	os.Setenv("DELIVERY_CHANNEL_SIZE", "10")
	os.Setenv("WORKER_FLUSH_TIMEOUT", "100")
	viper.AutomaticEnv()
	wc := WorkerConfigLoader()
	assert.Equal(t, 100, wc.WorkerFlushTimeout())
	assert.Equal(t, 10, wc.DeliveryChannelSize())
	assert.Equal(t, 5, wc.ChannelSize())
	assert.Equal(t, 2, wc.WorkersPoolSize())

	viper.Reset()
}
