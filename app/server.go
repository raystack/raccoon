package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"raccoon/config"
	"raccoon/logger"
	"raccoon/publisher"
	ws "raccoon/websocket"
	"raccoon/worker"
	"syscall"
	"time"
)

// StartServer starts the server
func StartServer(ctx context.Context, cancel context.CancelFunc) {

	//@TODO - create publisher with ctx

	//@TODO - create events-channels, workers (go routines) with ctx

	//start server @TODOD - the wss handler should be passed with the events-channel
	wssServer, bufferChannel := ws.CreateServer()
	logger.Info("Start Server -->")
	wssServer.StartHTTPServer(ctx, cancel)
	logger.Info("Start publisher -->")
	kPublisher, err := publisher.NewKafka(config.NewKafkaConfig())
	if err != nil {
		logger.Error("Error creating kafka producer", err)
		logger.Info("Exiting server")
		os.Exit(0)
	}

	logger.Info("Start worker -->")
	workerPool := worker.CreateWorkerPool(config.WorkerConfigLoader().WorkersPoolSize(), bufferChannel, config.WorkerConfigLoader().DeliveryChannelSize(), kPublisher)
	workerPool.StartWorkers()

	go shutDownServer(ctx, cancel, workerPool, kPublisher)
}

func shutDownServer(ctx context.Context, cancel context.CancelFunc, workerPool *worker.Pool, kp *publisher.Kafka) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info(fmt.Sprintf("[App.Server] Received a signal %s", sig))
			time.Sleep(3 * time.Second)
			// Temporary graceful shutdown mechanism
			workerPool.Flush()
			flushInterval := config.NewKafkaConfig().GetFlushInterval()
			logger.Info("Closing Kafka producer")
			logger.Info(fmt.Sprintf("Wait %d ms for all messages to be delivered", flushInterval))
			kp.Close()
			logger.Info("Exiting server")
			os.Exit(0)
		default:
			logger.Info(fmt.Sprintf("[App.Server] Received a unexpected signal %s", sig))
		}
	}
}
