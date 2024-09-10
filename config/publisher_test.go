package config

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
)

func TestKafkaPublisher(t *testing.T) {
	t.Run("should convert client config to equivalent config map", func(t *testing.T) {
		p := publisherKafka{
			ClientConfig: kafkaClientConfig{
				BootstrapServers: "localhost:8082",
				Acks:             "1",
			},
		}

		cm := p.ToKafkaConfigMap()
		expected := &kafka.ConfigMap{}
		expected.SetKey("bootstrap.servers", "localhost:8082")
		expected.SetKey("acks", "1")
		assert.Equal(t, cm, expected)
	})

	t.Run("should only add configs to config map that have non-zero values", func(t *testing.T) {
		assert.Equal(t, &kafka.ConfigMap{}, publisherKafka{}.ToKafkaConfigMap())
	})
}
