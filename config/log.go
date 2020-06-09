package config

func LogLevel() string {
	return mustGetString("LOG_LEVEL")
}
