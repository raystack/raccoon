package app

import (
	"context"

	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/pkg/logger"
)

// Run the server
func Run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	//@TODO - init config

	//start server
	logger.Infof("Raccoon %s", config.Version)
	StartServer(ctx, cancel)
	logger.Info("App.Run --> Complete")
	<-ctx.Done()
	return nil
}
