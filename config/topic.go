package config

import (
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

type TopicConfig struct {
	Format            string
	NumPartitions     int
	ReplicationFactor int
}

func (tc TopicConfig) GetFormat() string {
	return tc.Format
}

func (tc TopicConfig) GetNumPartitions() int {
	return tc.NumPartitions
}

func (tc TopicConfig) GetReplicationFactor() int {
	return tc.ReplicationFactor
}

func NewTopicConfig() TopicConfig {
	viper.SetDefault("TOPIC_FORMAT", "%s")
	tc := TopicConfig{
		Format:            mustGetString("TOPIC_FORMAT"),
		NumPartitions:     mustGetInt("TOPIC_NUM_PARTITIONS"),
		ReplicationFactor: mustGetInt("TOPIC_REPLICATION_FACTOR"),
	}

	return tc
}

func (tc TopicConfig) ToTopicConfigMap() map[string]string {
	cm := make(map[string]string)
	for key, val := range allSettings() {
		var v string
		if strings.HasPrefix(key, "topic_cm_") {
			switch a := val.(type) {
			case int:
				v = strconv.Itoa(a)
			case bool:
				v = strconv.FormatBool(a)
			case string:
				v = val.(string)
			default:
				panic("Type on TOPIC_CM_ is not supported yet")
			}
			cm[strings.Join(strings.Split(key, "_")[2:], ".")] = v
		}
	}
	return cm
}
