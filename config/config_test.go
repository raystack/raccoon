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
	os.Setenv("KAFKA_BROKER_LIST", "kafka:6668")
	expectedKafkaConfig := KafkaConfig{
		brokerList: "kafka:6668",
	}
	kafkaConfig := NewKafkaConfig()
	viper.AutomaticEnv()
	assert.Equal(t, expectedKafkaConfig, kafkaConfig)
}
