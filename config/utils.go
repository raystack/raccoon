package config

import (
	"os"
	"strings"
)

func dynamicConfigLoad(prefix string) []byte {
	var kafkaConfigs []string
	for _, v := range os.Environ() {
		if strings.HasPrefix(strings.ToLower(v), prefix) {
			kafkaConfigs = append(kafkaConfigs, v)
		}
	}
	yamlFormatted := []byte(
		strings.Replace(strings.Join(kafkaConfigs, "\n"), "=", ": ", -1))
	return yamlFormatted
}
