package config_test

import (
	"testing"

	"github.com/raystack/raccoon/config"
)

func TestWalk(t *testing.T) {
	var cfgSet = make(map[string]bool)
	for _, cfg := range config.Walk() {
		cfgSet[cfg.Meta.Tag.Get("cmdx")] = true
	}

	var samples = []string{
		"publisher.type",
		"server.websocket.conn.id.header",
		"publisher.kafka.client.bootstrap.servers",
		"publisher.pubsub.project.id",
		"publisher.pubsub.credentials",
		"publisher.kinesis.credentials",
	}

	for _, cfg := range samples {
		if !cfgSet[cfg] {
			t.Errorf("expected Walk() to return config with cmdx tag %s, but it didn't", cfg)
		}
	}
}
