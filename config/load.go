package config

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	defaults "github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

var loaded bool

// Load configs from env or yaml and set it to respective keys
func Load() error {
	if loaded {
		return nil
	}
	loaded = true
	viper.AutomaticEnv()
	viper.SetConfigName(".env")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.SetConfigType("env")
	viper.ReadInConfig()

	viper.MergeConfig(bytes.NewBuffer(dynamicKafkaClientConfigLoad()))

	return validate(&Server)
}

func init() {
	defaults.SetDefaults(&Server)
}

func validate(srv *server) error {
	if strings.TrimSpace(srv.Websocket.ConnIDHeader) == "" {
		return errFieldRequired(srv.Websocket, "ConnIDHeader")
	}
	if srv.Publisher == "pubsub" {
		if strings.TrimSpace(srv.PublisherPubSub.ProjectId) == "" {
			return errFieldRequired(srv.PublisherPubSub, "ProjectId")
		}
		if strings.TrimSpace(srv.PublisherPubSub.CredentialsFile) == "" {
			return errFieldRequired(srv.PublisherPubSub, "CredentialsFile")
		}
	}

	// there are no concrete fields that refer to this config
	kafkaServers := "PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS"
	if srv.Publisher == "kafka" && !viper.IsSet(kafkaServers) {
		flag := strings.ToLower(kafkaServers)
		flag = strings.ReplaceAll(flag, "_", ".")
		return errRequired(kafkaServers, flag)
	}

	return nil
}

func errFieldRequired(cfg any, field string) error {
	f, found := reflect.TypeOf(cfg).FieldByName(field)
	if !found {
		msg := fmt.Sprintf("unknown field %s in %s", field, cfg)
		panic(msg)
	}
	return fmt.Errorf("%s (--%s) is required", f.Tag.Get("mapstructure"), f.Tag.Get("cmdx"))
}

func errRequired(env, cmd string) error {
	return fmt.Errorf("%s (--%s) is required", env, cmd)
}
