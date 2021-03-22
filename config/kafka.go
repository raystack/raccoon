package config

import (
	"os"
	"raccoon/config/util"
	"strings"

	"github.com/spf13/viper"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type KafkaConfig struct {
	FlushInterval int
	TopicFormat   string
}

func (kc KafkaConfig) GetFlushInterval() int {
	return kc.FlushInterval
}

func (kc KafkaConfig) GetTopicFormat() string {
	return kc.TopicFormat
}

func NewKafkaConfig() KafkaConfig {
	viper.SetDefault("topic_format", "%s")
	kc := KafkaConfig{
		FlushInterval: util.MustGetInt("KAFKA_FLUSH_INTERVAL"),
		TopicFormat:   util.MustGetString("TOPIC_FORMAT"),
	}
	return kc
}

func (kc KafkaConfig) ToKafkaConfigMap() *kafka.ConfigMap {
	configMap := &kafka.ConfigMap{}
	for key, value := range viper.AllSettings() {
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
