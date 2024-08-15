package config

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/spf13/viper"
)

// configuration wrapper for initialising global configs
type cfg struct {
	Server    *server    `mapstructure:"server"`
	Publisher *publisher `mapstructure:"publisher"`
	Worker    *worker    `mapstructure:"worker"`
	Event     *event     `mapstructure:"event"`
	Metric    *metric    `mapstructure:"metric"`
	Log       *log       `mapstructure:"log"`
}

// prepare applies defaults and fallback values to global configurations
// prepare must be called after loading configs using viper
func prepare() {

	// parse kafka dynamic config
	viper.MergeConfig(bytes.NewBuffer(dynamicKafkaClientConfigLoad()))

	// add fallback for pubsub credentials
	if Publisher.Type == "pubsub" {
		defaultCredentials := strings.TrimSpace(
			os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		)
		if strings.TrimSpace(Publisher.PubSub.CredentialsFile) == "" && defaultCredentials != "" {
			Publisher.PubSub.CredentialsFile = defaultCredentials
		}
	}

	// add default CORS headers
	corsHeaders := []string{"Content-Type"}
	provisionalHeaders := []string{
		"SERVER_WEBSOCKET_CONN_GROUP_HEADER",
		"SERVER_WEBSOCKET_CONN_ID_HEADER",
	}

	for _, header := range provisionalHeaders {
		value := strings.TrimSpace(os.Getenv(header))
		if value != "" {
			corsHeaders = append(corsHeaders, value)
		}
	}

	for _, header := range corsHeaders {
		if !slices.Contains(Server.CORS.AllowedHeaders, header) {
			Server.CORS.AllowedHeaders = append(Server.CORS.AllowedHeaders, header)
		}
	}

}

// validate global configurations
func validate() error {
	trim := strings.TrimSpace
	if trim(Server.Websocket.Conn.IDHeader) == "" {
		return errCfgRequired("Server.Websocket.Conn.IDHeader")
	}
	if Publisher.Type == "pubsub" {
		if trim(Publisher.PubSub.ProjectId) == "" {
			return errCfgRequired("Publisher.PubSub.ProjectId")
		}
		if trim(Publisher.PubSub.CredentialsFile) == "" {
			return errCfgRequired("Publisher.PubSub.CredentialsFile")
		}
	}

	if Publisher.Type == "kinesis" {

		hasAWSEnvCreds := trim(os.Getenv("AWS_ACCESS_KEY_ID")) != "" &&
			trim(os.Getenv("AWS_SECRET_ACCESS_KEY")) != ""

		if trim(Publisher.Kinesis.CredentialsFile) == "" && !hasAWSEnvCreds {
			return errCfgRequired("Publisher.Kinesis.CredentialsFile")
		}
	}

	// there are no concrete fields that refer to this config
	kafkaServers := "PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS"
	if Publisher.Type == "kafka" && !viper.IsSet(kafkaServers) {
		flag := strings.ToLower(kafkaServers)
		flag = strings.ReplaceAll(flag, "_", ".")
		return errRequired(kafkaServers, flag)
	}

	validPublishers := []string{
		"kafka",
		"kinesis",
		"pubsub",
	}
	if !slices.Contains(validPublishers, Publisher.Type) {
		return fmt.Errorf("unknown publisher: %s", Publisher.Type)
	}

	return nil
}
