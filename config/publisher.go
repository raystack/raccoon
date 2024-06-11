package config

import (
	"bytes"
	"os"
	"strings"
	"time"

	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/raystack/raccoon/config/util"
	"github.com/spf13/viper"
)

var Publisher string
var PublisherKafka publisherKafka
var PublisherPubSub publisherPubSub
var dynamicKafkaClientConfigPrefix = "PUBLISHER_KAFKA_CLIENT_"

type publisherPubSub struct {
	ProjectId            string
	TopicAutoCreate      bool
	TopicRetentionPeriod time.Duration
}

type publisherKafka struct {
	FlushInterval int
}

func (k publisherKafka) ToKafkaConfigMap() *confluent.ConfigMap {
	configMap := &confluent.ConfigMap{}
	for key, value := range viper.AllSettings() {
		if strings.HasPrefix(strings.ToUpper(key), dynamicKafkaClientConfigPrefix) {
			clientConfig := key[len(dynamicKafkaClientConfigPrefix):]
			configMap.SetKey(strings.Join(strings.Split(clientConfig, "_"), "."), value)
		}
	}
	return configMap
}

func dynamicKafkaClientConfigLoad() []byte {
	var kafkaConfigs []string
	for _, v := range os.Environ() {
		if strings.HasPrefix(strings.ToUpper(v), dynamicKafkaClientConfigPrefix) {
			kafkaConfigs = append(kafkaConfigs, v)
		}
	}
	yamlFormatted := []byte(
		strings.Replace(strings.Join(kafkaConfigs, "\n"), "=", ": ", -1))
	return yamlFormatted
}

func publisherKafkaConfigLoader() {
	viper.SetDefault("PUBLISHER_KAFKA_CLIENT_QUEUE_BUFFERING_MAX_MESSAGES", "100000")
	viper.SetDefault("PUBLISHER_KAFKA_FLUSH_INTERVAL_MS", "1000")
	viper.MergeConfig(bytes.NewBuffer(dynamicKafkaClientConfigLoad()))

	PublisherKafka = publisherKafka{
		FlushInterval: util.MustGetInt("PUBLISHER_KAFKA_FLUSH_INTERVAL_MS"),
	}
}

func publisherPubSubLoader() {
	envTopicAutoCreate := "PUBLISHER_PUBSUB_TOPIC_AUTOCREATE"
	envTopicRetentionDuration := "PUBLISHER_PUBSUB_TOPIC_RETENTION_MS"

	viper.SetDefault(envTopicAutoCreate, "false")
	viper.SetDefault(envTopicRetentionDuration, "0")
	PublisherPubSub = publisherPubSub{
		ProjectId:            util.MustGetString("PUBLISHER_PUBSUB_PROJECT_ID"),
		TopicAutoCreate:      util.MustGetBool(envTopicAutoCreate),
		TopicRetentionPeriod: util.MustGetDuration(envTopicRetentionDuration, time.Millisecond),
	}
}

func publisherConfigLoader() {

	viper.SetDefault("PUBLISHER_TYPE", "kafka")

	Publisher = util.MustGetString("PUBLISHER_TYPE")
	Publisher = strings.ToLower(
		strings.TrimSpace(Publisher),
	)

	switch Publisher {
	case "kafka":
		publisherKafkaConfigLoader()
	case "pubsub":
		publisherPubSubLoader()
	}
}
