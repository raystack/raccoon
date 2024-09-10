package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCfgMetadata(t *testing.T) {
	t.Run("should return an error for a non-existent field", func(t *testing.T) {
		_, _, err := cfgMetadata("unknown.field")
		assert.Error(t, err)
	})
	t.Run("should return an error if the terminal field is missing cmdx tag", func(t *testing.T) {
		// this is an intermediate field, so it shouldn't contain cmdx tag
		_, _, err := cfgMetadata("Publisher.Kafka")
		assert.Error(t, err)
	})

}
