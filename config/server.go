package config

func AppPort() string {
	return mustGetString("APP_PORT")
}
