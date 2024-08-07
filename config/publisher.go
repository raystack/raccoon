package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/raystack/raccoon/config/util"
	"github.com/spf13/viper"
)

var Publisher string
var PublisherKafka publisherKafka
var PublisherPubSub publisherPubSub
var PublisherKinesis publisherKinesis
var dynamicKafkaClientConfigPrefix = "PUBLISHER_KAFKA_CLIENT_"

type publisherPubSub struct {
	ProjectId             string
	TopicAutoCreate       bool
	TopicRetentionPeriod  time.Duration
	PublishTimeout        time.Duration
	PublishDelayThreshold time.Duration
	PublishCountThreshold int
	PublishByteThreshold  int
	CredentialsFile       string
}

type publisherKinesis struct {
	Region              string
	CredentialsFile     string
	StreamAutoCreate    bool
	StreamProbeInterval time.Duration
	StreamMode          string
	DefaultShards       uint32
	PublishTimeout      time.Duration
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
	var (
		envCredentialsFile        = "PUBLISHER_PUBSUB_CREDENTIALS"
		envTopicAutoCreate        = "PUBLISHER_PUBSUB_TOPIC_AUTOCREATE"
		envTopicRetentionDuration = "PUBLISHER_PUBSUB_TOPIC_RETENTION_MS"
		envPublishDelayThreshold  = "PUBLISHER_PUBSUB_PUBLISH_DELAY_THRESHOLD_MS"
		envPublishCountThreshold  = "PUBLISHER_PUBSUB_PUBLISH_COUNT_THRESHOLD"
		envPublishByteThreshold   = "PUBLISHER_PUBSUB_PUBLISH_BYTE_THRESHOLD"
		envPublishTimeout         = "PUBLISHER_PUBSUB_PUBLISH_TIMEOUT_MS"
	)

	defaultCredentials := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if strings.TrimSpace(defaultCredentials) != "" {
		viper.SetDefault(envCredentialsFile, defaultCredentials)
	}

	viper.SetDefault(envTopicAutoCreate, "false")
	viper.SetDefault(envTopicRetentionDuration, "0")
	viper.SetDefault(envPublishDelayThreshold, "10")
	viper.SetDefault(envPublishCountThreshold, "100")
	viper.SetDefault(envPublishByteThreshold, "1000000") // ~1mb
	viper.SetDefault(envPublishTimeout, "60000")         // 1 minute

	PublisherPubSub = publisherPubSub{
		ProjectId:             util.MustGetString("PUBLISHER_PUBSUB_PROJECT_ID"),
		CredentialsFile:       util.MustGetString(envCredentialsFile),
		TopicAutoCreate:       util.MustGetBool(envTopicAutoCreate),
		TopicRetentionPeriod:  util.MustGetDuration(envTopicRetentionDuration, time.Millisecond),
		PublishTimeout:        util.MustGetDuration(envPublishTimeout, time.Millisecond),
		PublishDelayThreshold: util.MustGetDuration(envPublishDelayThreshold, time.Millisecond),
		PublishCountThreshold: util.MustGetInt(envPublishCountThreshold),
		PublishByteThreshold:  util.MustGetInt(envPublishByteThreshold),
	}
}

func publisherKinesisLoader() {
	var (
		envAWSRegion           = "PUBLISHER_KINESIS_AWS_REGION"
		envCredentialsFile     = "PUBLISHER_KINESIS_CREDENTIALS"
		envStreamProbeInterval = "PUBLISHER_KINESIS_STREAM_PROBE_INTERVAL_MS"
		envStreamAutoCreate    = "PUBLISHER_KINESIS_STREAM_AUTOCREATE"
		envStreamMode          = "PUBLISHER_KINESIS_STREAM_MODE"
		envStreamDefaultShards = "PUBLISHER_KINESIS_STREAM_SHARDS"
		envPublishTimeout      = "PUBLISHER_KINESIS_PUBLISH_TIMEOUT_MS"
	)

	defaultRegion := os.Getenv("AWS_REGION")
	if strings.TrimSpace(defaultRegion) != "" {
		viper.SetDefault(envAWSRegion, defaultRegion)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(
			fmt.Sprintf("unable to locate user home directory: %v", err),
		)
	}
	defaultCredentials := filepath.Join(home, ".aws", "credentials")

	viper.SetDefault(envCredentialsFile, defaultCredentials)
	viper.SetDefault(envStreamProbeInterval, "1000")
	viper.SetDefault(envStreamAutoCreate, "false")
	viper.SetDefault(envStreamMode, "ON_DEMAND")
	viper.SetDefault(envStreamDefaultShards, "4")
	viper.SetDefault(envPublishTimeout, "60000")

	PublisherKinesis = publisherKinesis{
		Region:              util.MustGetString(envAWSRegion),
		CredentialsFile:     util.MustGetString(envCredentialsFile),
		StreamAutoCreate:    util.MustGetBool(envStreamAutoCreate),
		StreamProbeInterval: util.MustGetDuration(envStreamProbeInterval, time.Millisecond),
		StreamMode:          util.MustGetString(envStreamMode),
		DefaultShards:       uint32(util.MustGetInt(envStreamDefaultShards)),
		PublishTimeout:      util.MustGetDuration(envPublishTimeout, time.Millisecond),
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
	case "kinesis":
		publisherKinesisLoader()
	}
}
