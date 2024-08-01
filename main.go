package main

import (
	"os"

	"github.com/raystack/raccoon/cmd"
)

// func main() {
// 	config.Load()
// 	middleware.Load()
// 	metrics.Setup()
// 	logger.SetLevel(config.Log.Level)
// 	err := app.Run()
// 	metrics.Close()
// 	if err != nil {
// 		logger.Fatal("init failure", err)
// 	}
// }

func main() {
	if err := cmd.New().Execute(); err != nil {
		os.Exit(1)
	}
}
