package util

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfigUtil(t *testing.T) {
	t.Run("MustGetString", func(t *testing.T) {
		t.Run("Should get correct value", func(t *testing.T) {
			os.Setenv("STR_CONFIG", "value string")
			viper.AutomaticEnv()
			assert.Equal(t, "value string", MustGetString("STR_CONFIG"))
		})
	})

	t.Run("MustGetInt", func(t *testing.T) {
		t.Run("Should get correct value", func(t *testing.T) {
			os.Setenv("INT_CONFIG", "3")
			viper.AutomaticEnv()
			assert.Equal(t, 3, MustGetInt("INT_CONFIG"))
		})
		t.Run("Should panic when value is not int", func(t *testing.T) {
			os.Setenv("INT_CONFIG", "value string")
			viper.AutomaticEnv()
			assert.Panics(t, func() { MustGetInt("INT_CONFIG") })
		})
	})

	t.Run("MustGetBool", func(t *testing.T) {
		t.Run("Should get correct value", func(t *testing.T) {
			os.Setenv("BOOL_CONFIG", "true")
			viper.AutomaticEnv()
			assert.Equal(t, true, MustGetBool("BOOL_CONFIG"))
		})
	})

	t.Run("MustGetDurationInSeconds", func(t *testing.T) {
		t.Run("Should get correct value", func(t *testing.T) {
			os.Setenv("DURATION_CONFIG", "20")
			viper.AutomaticEnv()
			assert.IsType(t, time.Second, MustGetDurationInSeconds("DURATION_CONFIG"))
		})
	})
}
