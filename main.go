package main

import (
	"raccoon/app"
	"raccoon/config"
	"raccoon/logger"
)

func main() {
	config.Load()
	logger.Setup()
	err := app.Run()
	if err != nil {
		logger.Fatal("init failure", err)
	}
}
