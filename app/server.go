package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"raccoon/config"
	raccoonhttp "raccoon/http"
	"raccoon/logger"
	"raccoon/metrics"
	"raccoon/pkg/collection"
	"raccoon/publisher"
	"raccoon/worker"
	"syscall"

	"google.golang.org/grpc"
)

// StartServer starts the server
func StartServer(ctx context.Context, cancel context.CancelFunc) {
	bufferChannel := make(chan *collection.EventsBatch)
	httpserver := raccoonhttp.CreateServer(bufferChannel)
	logger.Info("Start Server -->")
	httpserver.StartServers(ctx, cancel)
	logger.Info("Start publisher -->")
	kPublisher, err := publisher.NewKafka()
	if err != nil {
		logger.Error("Error creating kafka producer", err)
		logger.Info("Exiting server")
		os.Exit(0)
	}

	logger.Info("Start worker -->")
	workerPool := worker.CreateWorkerPool(config.Worker.WorkersPoolSize, bufferChannel, config.Worker.DeliveryChannelSize, kPublisher)
	workerPool.StartWorkers()
	go kPublisher.ReportStats()
	go shutDownServer(ctx, cancel, httpserver.HTTPServer, bufferChannel, workerPool, kPublisher, httpserver.GRPCServer)
}

func shutDownServer(ctx context.Context, cancel context.CancelFunc, httpServer *http.Server, bufferChannel chan *collection.EventsBatch, workerPool *worker.Pool, kp *publisher.Kafka, grpcServer *grpc.Server) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info(fmt.Sprintf("[App.Server] Received a signal %s", sig))
			httpServer.Shutdown(ctx)
			grpcServer.GracefulStop()
			logger.Info("Server shutdown all the listeners")
			timedOut := workerPool.FlushWithTimeOut(config.Worker.WorkerFlushTimeout)
			if timedOut {
				logger.Info(fmt.Sprintf("WorkerPool flush timedout %t", timedOut))
			}
			flushInterval := config.PublisherKafka.FlushInterval
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
