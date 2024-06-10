package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/raystack/raccoon/collection"
	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/logger"
	"github.com/raystack/raccoon/metrics"
	"github.com/raystack/raccoon/publisher"
	"github.com/raystack/raccoon/services"
	"github.com/raystack/raccoon/worker"
)

type Publisher interface {
	worker.Producer
	io.Closer
}

// StartServer starts the server
func StartServer(ctx context.Context, cancel context.CancelFunc) {
	bufferChannel := make(chan collection.CollectRequest, config.Worker.ChannelSize)
	httpServices := services.Create(bufferChannel)
	logger.Info("Start Server -->")
	httpServices.Start(ctx, cancel)
	logger.Infof("Start publisher --> %s", config.Publisher)
	publisher, err := initPublisher()
	if err != nil {
		logger.Errorf("Error creating %q publisher: %v\n", config.Publisher, err)
		logger.Info("Exiting server")
		os.Exit(0)
	}

	logger.Info("Start worker -->")
	workerPool := worker.CreateWorkerPool(config.Worker.WorkersPoolSize, bufferChannel, config.Worker.DeliveryChannelSize, publisher)
	workerPool.StartWorkers()

	go reportProcMetrics()
	go shutDownServer(ctx, cancel, httpServices, bufferChannel, workerPool, publisher)
}

func shutDownServer(ctx context.Context, cancel context.CancelFunc, httpServices services.Services, bufferChannel chan collection.CollectRequest, workerPool *worker.Pool, kp Publisher) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info(fmt.Sprintf("[App.Server] Received a signal %s", sig))
			httpServices.Shutdown(ctx)
			logger.Info("Server shutdown all the listeners")
			timedOut := workerPool.FlushWithTimeOut(config.Worker.WorkerFlushTimeout)
			if timedOut {
				logger.Info(fmt.Sprintf("WorkerPool flush timedout %t", timedOut))
			}
			flushInterval := config.PublisherKafka.FlushInterval
			logger.Info("Closing Kafka producer")
			logger.Info(fmt.Sprintf("Wait %d ms for all messages to be delivered", flushInterval))
			eventsInProducer := 0

			err := kp.Close()
			if err != nil {
				switch e := err.(type) {
				case *publisher.UnflushedEventsError:
					eventsInProducer = e.Count
				default:
					logger.Errorf("error closing %q publisher: %v", config.Publisher, err)
				}
			}

			/**
			@TODO - should compute the actual no., of events per batch and therefore the total. We can do this only when we close all the active connections
			Until then we fall back to approximation */
			eventsInChannel := len(bufferChannel) * 7
			logger.Info(fmt.Sprintf("Outstanding unprocessed events in the channel, data lost ~ (No batches %d * 5 events) = ~%d", len(bufferChannel), eventsInChannel))
			metrics.Count("kafka_messages_delivered_total", int64(eventsInChannel+eventsInProducer), map[string]string{"success": "false", "conn_group": "NA", "event_type": "NA"})
			logger.Info("Exiting server")
			cancel()
		default:
			logger.Info(fmt.Sprintf("[App.Server] Received a unexpected signal %s", sig))
		}
	}
}

func reportProcMetrics() {
	m := &runtime.MemStats{}
	for range time.Tick(config.MetricInfo.RuntimeStatsRecordInterval) {
		metrics.Gauge("server_go_routines_count_current", runtime.NumGoroutine(), map[string]string{})
		runtime.ReadMemStats(m)
		metrics.Gauge("server_mem_heap_alloc_bytes_current", m.HeapAlloc, map[string]string{})
		metrics.Gauge("server_mem_heap_inuse_bytes_current", m.HeapInuse, map[string]string{})
		metrics.Gauge("server_mem_heap_objects_total_current", m.HeapObjects, map[string]string{})
		metrics.Gauge("server_mem_stack_inuse_bytes_current", m.StackInuse, map[string]string{})
		metrics.Gauge("server_mem_gc_triggered_current", m.LastGC/1000, map[string]string{})
		metrics.Gauge("server_mem_gc_pauseNs_current", m.PauseNs[(m.NumGC+255)%256]/1000, map[string]string{})
		metrics.Gauge("server_mem_gc_count_current", m.NumGC, map[string]string{})
		metrics.Gauge("server_mem_gc_pauseTotalNs_current", m.PauseTotalNs, map[string]string{})
	}
}

func initPublisher() (Publisher, error) {
	switch config.Publisher {
	case "kafka":
		return publisher.NewKafka()
	case "pubsub":
		return publisher.NewPubSub(
			config.PublisherPubSub.ProjectId,
			config.EventDistribution.PublisherPattern,
		)
	default:
		return nil, fmt.Errorf("unknown publisher: %v", config.Publisher)
	}
}
