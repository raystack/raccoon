package config

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLogLevel(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	viper.AutomaticEnv()
	assert.Equal(t, "debug", LogLevel())

	viper.Reset()
}

func TestAppPort(t *testing.T) {
	os.Setenv("APP_PORT", "8080")
	viper.AutomaticEnv()
	assert.Equal(t, "8080", ServerConfigLoader().AppPort)

	viper.Reset()
}

func TestNewKafkaConfig(t *testing.T) {
	os.Setenv("KAFKA_TOPIC", "test1")
	os.Setenv("KAFKA_FLUSH_INTERVAL", "1000")

	expectedKafkaConfig := KafkaConfig{
		Topic:         "test1",
		FlushInterval: 1000,
	}

	viper.AutomaticEnv()
	kafkaConfig := NewKafkaConfig()
	assert.Equal(t, expectedKafkaConfig, kafkaConfig)

	viper.Reset()
}

func TestDynamicKafkaConfigLoad(t *testing.T) {
	os.Setenv("KAFKA_CLIENT_RANDOM", "anything")
	os.Setenv("KAFKA_CLIENT_BOOTSTRAP_SERVERS", "localhost:9092")
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(dynamicKafkaConfigLoad()))
	assert.Equal(t, viper.AllSettings()["kafka_client_random"], "anything")
	assert.Equal(t, viper.AllSettings()["kafka_client_bootstrap_servers"], "localhost:9092")

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
