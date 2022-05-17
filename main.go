package main

import (
	"github.com/odpf/raccoon/app"
	"github.com/odpf/raccoon/config"
	"github.com/odpf/raccoon/logger"
	"github.com/odpf/raccoon/metrics"
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
