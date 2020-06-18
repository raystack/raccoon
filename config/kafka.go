package config

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"strings"
)

type KafkaConfig struct {
	topic         string
	flushInterval int
}

func (kc KafkaConfig) Topic() string {
	return kc.topic
}

func (kc KafkaConfig) FlushInterval() int {
	return kc.flushInterval
}

func NewKafkaConfig() KafkaConfig {
	kc := KafkaConfig{
		topic:         mustGetString("KAFKA_TOPIC"),
		flushInterval: mustGetInt("KAFKA_FLUSH_INTERVAL"),
	}
	return kc
}

func (cfg KafkaConfig) ToKafkaConfigMap() *kafka.ConfigMap {
	configMap := &kafka.ConfigMap{}
	for key, value := range allSettings() {
		if len(key) > 13 && key[0:13] == "kafka_client_" {
			configMap.SetKey(strings.Join(strings.Split(key, "_")[2:], "."), value)
		}
	}
	return configMap
}
