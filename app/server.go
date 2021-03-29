package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"raccoon/config"
	"raccoon/logger"
	"raccoon/metrics"
	"raccoon/publisher"
	ws "raccoon/websocket"
	"raccoon/worker"
	"syscall"
	"time"
)

// StartServer starts the server
func StartServer(ctx context.Context, cancel context.CancelFunc) {
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
	go kPublisher.ReportStats()
	go shutDownServer(ctx, cancel, wssServer.HTTPServer, bufferChannel, workerPool, kPublisher)
}

func shutDownServer(ctx context.Context, cancel context.CancelFunc, wssServer *http.Server, bufferChannel chan ws.EventsBatch, workerPool *worker.Pool, kp *publisher.Kafka) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info(fmt.Sprintf("[App.Server] Received a signal %s", sig))
			wssServer.Shutdown(ctx)
			logger.Info("Server shutdown all the listeners")
			timedOut := workerPool.FlushWithTimeOut(time.Duration(config.WorkerConfigLoader().WorkerFlushTimeout()) * time.Second)
			if timedOut {
				logger.Info(fmt.Sprintf("WorkerPool flush timedout %t", timedOut))
			}
			flushInterval := config.NewKafkaConfig().GetFlushInterval()
			logger.Info("Closing Kafka producer")
			logger.Info(fmt.Sprintf("Wait %d ms for all messages to be delivered", flushInterval))
			eventsInProducer := kp.Close()
			/**
			@TODO - should compute the actual no., of events per batch and therefore the total. We can do this only when we close all the active connections
			Until then we fall back to approximation */
			eventsInChannel := len(bufferChannel) * 7
			logger.Info(fmt.Sprintf("Outstanding unprocessed events in the channel, data lost ~ (No batches %d * 5 events) = ~%d", len(bufferChannel), eventsInChannel))
			metrics.Count("kafka_messages_delivered_total", eventsInChannel+eventsInProducer, "success=false")
			logger.Info("Exiting server")
			os.Exit(0)
		default:
			logger.Info(fmt.Sprintf("[App.Server] Received a unexpected signal %s", sig))
		}
	}
}
