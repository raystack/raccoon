package main

import (
	"raccoon/app"
	_ "raccoon/config"
	"raccoon/logger"
	"raccoon/metrics"
)

func main() {
	metrics.Setup()
	err := app.Run()
	metrics.Close()
	if err != nil {
		logger.Fatal("init failure", err)
	}
}
