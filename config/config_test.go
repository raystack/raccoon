package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLogLevel(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	viper.AutomaticEnv()
	assert.Equal(t, "debug", LogLevel())
}

func TestAppPort(t *testing.T) {
	os.Setenv("APP_PORT", "8080")
	viper.AutomaticEnv()
	assert.Equal(t, "8080", AppPort())
}

func TestNewKafkaConfig(t *testing.T) {
	os.Setenv("KAFKA_TOPIC", "test1")
	os.Setenv("KAFKA_FLUSH_INTERVAL", "1000")

	expectedKafkaConfig := KafkaConfig{
		topic:         "test1",
		flushInterval: 1000,
	}

	viper.AutomaticEnv()
	kafkaConfig := NewKafkaConfig()
	assert.Equal(t, expectedKafkaConfig, kafkaConfig)
}

func TestKafkaConfig_ToKafkaConfigMap(t *testing.T) {
	os.Setenv("KAFKA_TOPIC", "test1")
	os.Setenv("KAFKA_FLUSH_INTERVAL", "1000")
	os.Setenv("KAFKA_CLIENT_BOOTSTRAP_SERVERS", "kafka:6668")
	os.Setenv("KAFKA_CLIENT_ACKS", "1")
	os.Setenv("KAFKA_CLIENT_QUEUE_BUFFERING_MAX_MESSAGES", "10000")

	viper.AutomaticEnv()
	viper.BindEnv("KAFKA_TOPIC")
	viper.BindEnv("KAFKA_FLUSH_INTERVAL")
	viper.BindEnv("KAFKA_CLIENT_BOOTSTRAP_SERVERS")
	viper.BindEnv("KAFKA_CLIENT_ACKS")
	viper.BindEnv("KAFKA_CLIENT_QUEUE_BUFFERING_MAX_MESSAGES")

	kafkaConfig := NewKafkaConfig().ToKafkaConfigMap()
	bootstrapServer, _ := kafkaConfig.Get("bootstrap.servers", "")
	topic, _ := kafkaConfig.Get("topic", "")
	assert.Equal(t, "kafka:6668", bootstrapServer)
	assert.Equal(t, "", topic)
	assert.Equal(t, 3, len(*kafkaConfig))
}
