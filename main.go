package main

import (
	"github.com/raystack/raccoon/app"
	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/logger"
	"github.com/raystack/raccoon/metrics"
)

func main() {
	config.Load()
	metrics.Setup()
	logger.SetLevel(config.Log.Level)
	err := app.Run()
	metrics.Close()
	if err != nil {
		logger.Fatal("init failure", err)
	}
}
