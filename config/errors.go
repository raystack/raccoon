package config

import (
	"fmt"
	"reflect"
)

func errFieldRequired(cfg any, field string) error {
	f, found := reflect.TypeOf(cfg).FieldByName(field)
	if !found {
		return fmt.Errorf("unknown field %s in %s", field, cfg)
	}
	requiredTags := []string{
		"mapstructure",
		"cmdx",
	}
	for _, tag := range requiredTags {
		if _, ok := f.Tag.Lookup(tag); !ok {
			return fmt.Errorf("%s.%s is missing tag %s", cfg, field, tag)
		}
	}
	return errRequired(f.Tag.Get("mapstructure"), f.Tag.Get("cmdx"))
}

func errRequired(env, cmd string) error {
	return fmt.Errorf("%s (--%s) is required", env, cmd)
}
