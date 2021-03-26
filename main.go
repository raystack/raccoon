package main

import (
	"raccoon/app"
	"raccoon/config"
	"raccoon/logger"
	"raccoon/metrics"
)

func main() {
	config.Load()
	metrics.Setup()
	err := app.Run()
	metrics.Close()
	if err != nil {
		logger.Fatal("init failure", err)
	}
}
