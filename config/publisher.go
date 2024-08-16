package config

import (
	"os"
	"strings"

	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/viper"
)

var Publisher publisher

var dynamicKafkaClientConfigPrefix = "PUBLISHER_KAFKA_CLIENT_"

type publisherPubSub struct {
	ProjectId               string `mapstructure:"project_id" cmdx:"publisher.pubsub.project.id"`
	TopicAutoCreate         bool   `mapstructure:"topic_autocreate" cmdx:"publisher.pubsub.topic.autocreate" default:"false"`
	TopicRetentionPeriodMS  int64  `mapstructure:"topic_retention_ms" cmdx:"publisher.pubsub.topic.retention.ms" default:"0"`
	PublishTimeoutMS        int64  `mapstructure:"publish_timeout_ms" cmdx:"publisher.pubsub.publish.timeout.ms" default:"60000"`
	PublishDelayThresholdMS int64  `mapstructure:"publish_delay_threshold_ms" cmdx:"publisher.pubsub.publish.delay.threshold.ms" default:"10"`
	PublishCountThreshold   int    `mapstructure:"publish_count_threshold" cmdx:"publisher.pubsub.publish.count.threshold" default:"100"`
	PublishByteThreshold    int    `mapstructure:"publish_byte_threshold" cmdx:"publisher.pubsub.publish.byte.threshold" default:"1000000"`
	CredentialsFile         string `mapstructure:"credentials" cmdx:"publisher.pubsub.credentials"`
}

type publisherKinesis struct {
	Region                string `mapstructure:"aws_region" cmdx:"publisher.kinesis.aws.region"`
	CredentialsFile       string `mapstructure:"credentials" cmdx:"publisher.kinesis.credentials"`
	StreamAutoCreate      bool   `mapstructure:"stream_autocreate" cmdx:"publisher.kinesis.stream.autocreate" default:"false"`
	StreamProbeIntervalMS int64  `mapstructure:"stream_probe_interval_ms" cmdx:"publisher.kinesis.stream.probe.interval.ms" default:"1000"`
	StreamMode            string `mapstructure:"stream_mode" cmdx:"publisher.kinesis.stream.mode" default:"ON_DEMAND"`
	DefaultShards         uint32 `mapstructure:"stream_shards" cmdx:"publisher.kinesis.stream.shards" default:"4"`
	PublishTimeoutMS      int64  `mapstructure:"publish_timeout_ms" cmdx:"publisher.kinesis.publish.timeout.ms" default:"60000"`
}

type publisherKafka struct {
	FlushInterval       int `mapstructure:"flush_interval_ms" cmdx:"publisher.kafka.flush.interval.ms" default:"1000"`
	DeliveryChannelSize int `mapstructure:"delivery_channel_size" cmdx:"publisher.kafka.delivery.channel.size" default:"10"`
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

type publisher struct {
	Type    string           `mapstructure:"type" cmdx:"publisher.type" default:"kafka"`
	Kafka   publisherKafka   `mapstructure:"kafka"`
	PubSub  publisherPubSub  `mapstructure:"pubsub"`
	Kinesis publisherKinesis `mapstructure:"kinesis"`
}
