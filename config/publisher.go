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
	ProjectId             string        `mapstructure:"PUBLISHER_PUBSUB_PROJECT_ID" cmdx:"publisher.pubsub.project.id"`
	TopicAutoCreate       bool          `mapstructure:"PUBLISHER_PUBSUB_TOPIC_AUTOCREATE" cmdx:"publisher.pubsub.topic.autocreate" default:"false"`
	TopicRetentionPeriod  time.Duration `mapstructure:"PUBLISHER_PUBSUB_TOPIC_RETENTION_MS" cmdx:"publisher.pubsub.topic.retention.ms" default:"0"`
	PublishTimeout        time.Duration `mapstructure:"PUBLISHER_PUBSUB_PUBLISH_TIMEOUT_MS" cmdx:"publisher.pubsub.publish.timeout.ms" default:"60000"`
	PublishDelayThreshold time.Duration `mapstructure:"PUBLISHER_PUBSUB_PUBLISH_DELAY_THRESHOLD_MS" cmdx:"publisher.pubsub.publish.delay.threshold.ms" default:"10"`
	PublishCountThreshold int           `mapstructure:"PUBLISHER_PUBSUB_PUBLISH_COUNT_THRESHOLD" cmdx:"publisher.pubsub.publish.count.threshold" default:"100"`
	PublishByteThreshold  int           `mapstructure:"PUBLISHER_PUBSUB_PUBLISH_BYTE_THRESHOLD" cmdx:"publisher.pubsub.publish.byte.threshold" default:"1000000"`
	CredentialsFile       string        `mapstructure:"PUBLISHER_PUBSUB_CREDENTIALS" cmdx:"publisher.pubsub.credentials"`
}

type publisherKinesis struct {
	Region              string        `mapstructure:"PUBLISHER_KINESIS_AWS_REGION" cmdx:"publisher.kinesis.aws.region"`
	CredentialsFile     string        `mapstructure:"PUBLISHER_KINESIS_CREDENTIALS" cmdx:"publisher.kinesis.credentials"`
	StreamAutoCreate    bool          `mapstructure:"PUBLISHER_KINESIS_STREAM_AUTOCREATE" cmdx:"publisher.kinesis.stream.autocreate" default:"false"`
	StreamProbeInterval time.Duration `mapstructure:"PUBLISHER_KINESIS_STREAM_PROBE_INTERVAL_MS" cmdx:"publisher.kinesis.stream.probe.interval.ms" default:"1000"`
	StreamMode          string        `mapstructure:"PUBLISHER_KINESIS_STREAM_MODE" cmdx:"publisher.kinesis.stream.mode" default:"ON_DEMAND"`
	DefaultShards       uint32        `mapstructure:"PUBLISHER_KINESIS_STREAM_SHARDS" cmdx:"publisher.kinesis.stream.shards" default:"4"`
	PublishTimeout      time.Duration `mapstructure:"PUBLISHER_KINESIS_PUBLISH_TIMEOUT_MS" cmdx:"publisher.kinesis.publish.timeout.ms" default:"60000"`
}

type publisherKafka struct {
	FlushInterval int `mapstructure:"PUBLISHER_KAFKA_FLUSH_INTERVAL_MS" cmdx:"publisher.kafka.flush.interval.ms" default:"1000"`
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
	viper.SetDefault("PUBLISHER_KAFKA_CLIENT_RETRY_BACKOFF_MS", "100")
	viper.SetDefault("PUBLISHER_KAFKA_CLIENT_ACKS", "-1")
	viper.SetDefault("PUBLISHER_KAFKA_CLIENT_RETRIES", "2147483647")
	viper.SetDefault("PUBLISHER_KAFKA_CLIENT_STATISTICS_INTERVAL_MS", "0")

	viper.MergeConfig(bytes.NewBuffer(dynamicKafkaClientConfigLoad()))

	// required
	_ = util.MustGetString("PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS")

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
