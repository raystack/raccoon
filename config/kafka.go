package config

import (
	"os"
	"strings"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type KafkaConfig struct {
	Topic         string
	FlushInterval int
}

func (kc KafkaConfig) GetTopic() string {
	return kc.Topic
}

func (kc KafkaConfig) GetFlushInterval() int {
	return kc.FlushInterval
}

func NewKafkaConfig() KafkaConfig {
	kc := KafkaConfig{
		Topic:         mustGetString("KAFKA_TOPIC"),
		FlushInterval: mustGetInt("KAFKA_FLUSH_INTERVAL"),
	}
	return kc
}

func (kc KafkaConfig) ToKafkaConfigMap() *kafka.ConfigMap {
	configMap := &kafka.ConfigMap{}
	for key, value := range allSettings() {
		if len(key) > 13 && key[0:13] == "kafka_client_" {
			configMap.SetKey(strings.Join(strings.Split(key, "_")[2:], "."), value)
		}
	}
	return configMap
}

func dynamicKafkaConfigLoad() []byte {
	var kafkaConfigs []string
	for _, v := range os.Environ() {
		if strings.HasPrefix(strings.ToLower(v), "kafka_client_") {
			kafkaConfigs = append(kafkaConfigs, v)
		}
	}
	yamlFormatted := []byte(
		strings.Replace(strings.Join(kafkaConfigs, "\n"), "=", ": ", -1))
	return yamlFormatted
}
