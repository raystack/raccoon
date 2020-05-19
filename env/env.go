package env

import (
	"github.com/pkg/errors"
	"os"
	"strconv"
)

func AppPort() int {
	return intValOrDefault(os.Getenv("APP_PORT"), 3000)
}

func stringValOrDefault(val string, def string) string {
	if val == "" {
		return def
	}
	return val
}

func intValOrDefault(val string, def int) int {
	if val == "" {
		return def
	}
	result, err := strconv.Atoi(val)
	if err != nil {
		panic(errors.Wrap(err, "expected int value for env value: "+val))
	}
	return result
}

func boolValOrDefault(val string, def bool) bool {
	if val == "" {
		return def
	}
	result, err := strconv.ParseBool(val)
	if err != nil {
		panic(errors.Wrap(err, "expected bool value for env value: "+val))
	}
	return result
}
