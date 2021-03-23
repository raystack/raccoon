package config

import (
	"bytes"
	"os"
	"raccoon/config/util"
	"strings"

	"github.com/spf13/viper"

	confluent "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var Kafka kafka

type kafka struct {
	FlushInterval int
}

func (k kafka) ToKafkaConfigMap() *confluent.ConfigMap {
	configMap := &confluent.ConfigMap{}
	for key, value := range viper.AllSettings() {
		if len(key) > 13 && key[0:13] == "kafka_client_" {
			configMap.SetKey(strings.Join(strings.Split(key, "_")[2:], "."), value)
		}
	}
	return configMap
}

func dynamicKafkaClientConfigLoad() []byte {
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

func kafkaConfigLoader() {
	viper.SetDefault("kafka_client_queue_buffering_max_messages", "100000")
	viper.SetDefault("kafka_flush_interval", "1000")
	viper.SetDefault("delivery_channel_size", "100")
	viper.MergeConfig(bytes.NewBuffer(dynamicKafkaClientConfigLoad()))

	Kafka = kafka{
		FlushInterval: util.MustGetInt("KAFKA_FLUSH_INTERVAL"),
	}
}
