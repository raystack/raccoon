package config

import (
	"fmt"
	"reflect"
	"strings"
)

func cfgMetadata(path string) (string, string, error) {

	var (
		fields = strings.Split(path, ".")
		parent = reflect.TypeOf(cfg)
		field  reflect.StructField
		found  bool
		env    []string
		hist   []string
		cmd    string
	)

	for len(fields) > 0 {
		fieldName := fields[0]
		hist = append(hist, fieldName)

		if field.Name == "" {
			field, found = parent.FieldByName(fieldName)
		} else {
			switch field.Type.Kind() {
			case reflect.Ptr:
				field, found = field.Type.Elem().FieldByName(fieldName)
			default:
				field, found = field.Type.FieldByName(fieldName)
			}
		}

		if !found {
			return "", "", fmt.Errorf("config is missing field %s", strings.Join(hist, "."))
		}

		envPartial := field.Tag.Get("mapstructure")
		if strings.TrimSpace(envPartial) == "" {
			return "", "", fmt.Errorf("config.%s is missing mapstructure tag or is empty", strings.Join(hist, "."))
		}
		env = append(env, strings.ToUpper(envPartial))

		if len(fields) == 1 {
			cmdxTag := field.Tag.Get("cmdx")
			if strings.TrimSpace(cmdxTag) == "" {
				return "", "", fmt.Errorf("config.%s is missing cmdx tag or is empty", strings.Join(hist, "."))
			}
			cmd = cmdxTag
		}
		fields = fields[1:]
	}
	return strings.Join(env, "_"), cmd, nil
}

func errCfgRequired(path string) error {
	env, cmd, err := cfgMetadata(path)
	if err != nil {
		return err
	}
	return ConfigError{env, cmd}
}

type ConfigError struct {
	Env, Flag string
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("%s (--%s) is required", e.Env, e.Flag)
}
