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

	"github.com/raystack/raccoon/collector"
	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/logger"
	"github.com/raystack/raccoon/metrics"
	"github.com/raystack/raccoon/publisher"
	"github.com/raystack/raccoon/publisher/kafka"
	"github.com/raystack/raccoon/publisher/kinesis"
	logpub "github.com/raystack/raccoon/publisher/log"
	"github.com/raystack/raccoon/publisher/pubsub"
	"github.com/raystack/raccoon/services"
	"github.com/raystack/raccoon/worker"

	pubsubsdk "cloud.google.com/go/pubsub"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	kinesissdk "github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"google.golang.org/api/option"
)

type Publisher interface {
	worker.Producer
	io.Closer
}

// StartServer starts the server
func StartServer(ctx context.Context, cancel context.CancelFunc) {
	bufferChannel := make(chan collector.CollectRequest, config.Worker.Buffer.ChannelSize)
	httpServices := services.Create(bufferChannel)
	logger.Info("Start Server -->")
	httpServices.Start(ctx, cancel)
	logger.Infof("Start publisher --> %s", config.Publisher.Type)
	publisher, err := initPublisher()
	if err != nil {
		logger.Errorf("Error creating %q publisher: %v\n", config.Publisher.Type, err)
		logger.Info("Exiting server")
		os.Exit(0)
	}

	logger.Info("Start worker -->")
	workerPool := worker.CreateWorkerPool(config.Worker.PoolSize, bufferChannel, publisher)
	workerPool.StartWorkers()

	go reportProcMetrics()
	go shutDownServer(ctx, cancel, httpServices, bufferChannel, workerPool, publisher)
}

func shutDownServer(ctx context.Context, cancel context.CancelFunc, httpServices services.Services, bufferChannel chan collector.CollectRequest, workerPool *worker.Pool, pub Publisher) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	workerFlushTimeout := time.Duration(config.Worker.Buffer.FlushTimeoutMS) * time.Millisecond
	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info(fmt.Sprintf("[App.Server] Received a signal %s", sig))
			httpServices.Shutdown(ctx)
			logger.Info("Server shutdown all the listeners")
			timedOut := workerPool.FlushWithTimeOut(workerFlushTimeout)
			if timedOut {
				logger.Info(fmt.Sprintf("WorkerPool flush timedout %t", timedOut))
			}
			flushInterval := config.Publisher.Kafka.FlushInterval
			logger.Infof("Closing %q producer\n", pub.Name())
			logger.Info(fmt.Sprintf("Wait %d ms for all messages to be delivered", flushInterval))
			eventsInProducer := 0

			err := pub.Close()
			if err != nil {
				switch e := err.(type) {
				case *publisher.UnflushedEventsError:
					eventsInProducer = e.Count
				default:
					logger.Errorf("error closing %q publisher: %v", pub.Name(), err)
				}
			}

			/**
			@TODO - should compute the actual no., of events per batch and therefore the total. We can do this only when we close all the active connections
			Until then we fall back to approximation */
			eventsInChannel := len(bufferChannel) * 7
			logger.Info(fmt.Sprintf("Outstanding unprocessed events in the channel, data lost ~ (No batches %d * 5 events) = ~%d", len(bufferChannel), eventsInChannel))
			metrics.Count(
				fmt.Sprintf("%s_messages_delivered_total", pub.Name()),
				int64(eventsInChannel+eventsInProducer),
				map[string]string{
					"conn_group": "NA",
					"event_type": "NA",
					"topic":      "NA",
				},
			)
			logger.Info("Exiting server")
			cancel()
		default:
			logger.Info(fmt.Sprintf("[App.Server] Received a unexpected signal %s", sig))
		}
	}
}

func reportProcMetrics() {
	m := &runtime.MemStats{}
	reportInterval := time.Duration(config.Metric.RuntimeStatsRecordIntervalMS) * time.Millisecond
	for range time.Tick(reportInterval) {
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
	switch config.Publisher.Type {
	case "kafka":
		return kafka.New()
	case "pubsub":
		client, err := pubsubsdk.NewClient(
			context.Background(),
			config.Publisher.PubSub.ProjectId,
			option.WithCredentialsFile(config.Publisher.PubSub.CredentialsFile),
		)
		if err != nil {
			return nil, fmt.Errorf("error creating pubsub client: %w", err)
		}
		var (
			topicRetention = time.Duration(config.Publisher.PubSub.TopicRetentionPeriodMS) * time.Millisecond
			delayThreshold = time.Duration(config.Publisher.PubSub.PublishDelayThresholdMS) * time.Millisecond
			publishTimeout = time.Duration(config.Publisher.PubSub.PublishTimeoutMS) * time.Millisecond
		)
		return pubsub.New(
			client,
			pubsub.WithTopicFormat(config.Event.DistributionPublisherPattern),
			pubsub.WithTopicAutocreate(config.Publisher.PubSub.TopicAutoCreate),
			pubsub.WithTopicRetention(topicRetention),
			pubsub.WithDelayThreshold(delayThreshold),
			pubsub.WithCountThreshold(config.Publisher.PubSub.PublishCountThreshold),
			pubsub.WithByteThreshold(config.Publisher.PubSub.PublishByteThreshold),
			pubsub.WithTimeout(publishTimeout),
		)
	case "kinesis":
		cfg, err := awsconfig.LoadDefaultConfig(
			context.Background(),
			awsconfig.WithRegion(config.Publisher.Kinesis.Region),
			awsconfig.WithSharedConfigFiles(
				[]string{config.Publisher.Kinesis.CredentialsFile},
			),
		)
		if err != nil {
			return nil, fmt.Errorf("error locating aws credentials: %w", err)
		}
		var (
			conf           = config.Publisher.Kinesis
			publishTimeout = time.Duration(conf.PublishTimeoutMS) * time.Millisecond
			probeInterval  = time.Duration(conf.StreamProbeIntervalMS) * time.Millisecond
		)

		return kinesis.New(
			kinesissdk.NewFromConfig(cfg),
			kinesis.WithStreamPattern(config.Event.DistributionPublisherPattern),
			kinesis.WithStreamAutocreate(conf.StreamAutoCreate),
			kinesis.WithStreamMode(types.StreamMode(conf.StreamMode)),
			kinesis.WithShards(conf.DefaultShards),
			kinesis.WithPublishTimeout(publishTimeout),
			kinesis.WithStreamProbleInterval(probeInterval),
		)
	case "log":
		return logpub.New(), nil
	default:
		return nil, fmt.Errorf("unknown publisher: %v", config.Publisher.Type)
	}
}
