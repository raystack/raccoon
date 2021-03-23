package util

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

func MustGetString(key string) string {
	mustHave(key)
	return viper.GetString(key)
}

func MustGetInt(key string) int {
	mustHave(key)
	v, err := strconv.Atoi(viper.GetString(key))
	if err != nil {
		panic(fmt.Sprintf("key %s is not a valid Integer value", key))
	}

	return v
}

func MustGetBool(key string) bool {
	mustHave(key)
	return viper.GetBool(key)
}

func MustGetDuration(key string, d time.Duration) time.Duration {
	return d * time.Duration(MustGetInt(key))
}

func mustHave(key string) {
	if !viper.IsSet(key) {
		panic(fmt.Sprintf("key %s is not set", key))
	}
}
