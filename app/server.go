package app

import (
	"clickstream-service/config"
	"clickstream-service/publisher"
	"clickstream-service/logger"
	ws "clickstream-service/websocket"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// StartServer starts the server
func StartServer(ctx context.Context, cancel context.CancelFunc) {
	
	//@TODO - create publisher with ctx

	//@TODO - create events-channels, workers (go routines) with ctx

	//start server @TODOD - the wss handler should be passed with the events-channel
	wssServer := ws.CreateServer()
	logger.Info("Start Server -->")
	wssServer.StartHTTPServer(ctx, cancel)
	logger.Info("Start publisher -->")

	publisher, err := publisher.NewProducer(ctx, config.NewKafkaConfig())
	if err != nil {
		logger.Error("Error while creating new producer", err)
	}

	publisher.PublishMessage([]byte("myvalue"), []byte("mykey"))
	go shutDownServer(ctx, cancel)
}

func shutDownServer(ctx context.Context, cancel context.CancelFunc) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info(fmt.Sprintf("[App.Server] Received a signal %s", sig))
			logger.Info(fmt.Sprintf("[App.Server] waiting for 3 secs grace period before shutdown ", ))
			time.Sleep(3 * time.Second)
			logger.Info("Exiting server")
			os.Exit(0)
		default:
			logger.Info(fmt.Sprintf("[App.Server] Received a unexpected signal %s", sig))
		}
	}
}
