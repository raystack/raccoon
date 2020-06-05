package main

import (
	"clickstream-service/app"
	"clickstream-service/config"
	"clickstream-service/logger"
)

func main() {
	config.Load()
	err := app.Run()
	if err != nil {
		logger.Fatal("init failure", err)
	}
}
