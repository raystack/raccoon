package config

import (
	"time"

	"github.com/raystack/raccoon/config/util"
	"github.com/spf13/viper"
)

// Worker contains configs for kafka publisher worker pool
var Worker worker

type worker struct {
	// WorkersPoolSize number of worker to push to kafka initiated at the start of Raccoon
	WorkersPoolSize int `mapstructure:"WORKER_POOL_SIZE" cmdx:"worker.pool.size" default:"5"`
	// ChannelSize channel size to buffer events before processed by worker
	ChannelSize int `mapstructure:"WORKER_BUFFER_CHANNEL_SIZE" cmdx:"worker.buffer.channel.size" default:"100"`
	//DeliveryChannelSize fetches the size of the delivery channel as configured
	DeliveryChannelSize int `mapstructure:"WORKER_KAFKA_DELIVERY_CHANNEL_SIZE" cmdx:"worker.kafka.delivery.channel.size" default:"10"`
	//WorkerFlushTimeout specifies a timeout interval that the workers use to timeout
	//in case the workers could not complete the flush. This enables a non-blocking flush.
	WorkerFlushTimeout time.Duration `mapstructure:"WORKER_BUFFER_FLUSH_TIMEOUT_MS" cmdx:"worker.buffer.flush.timeout.ms" default:"5000"`
}

// workerConfigLoader constructs a singleton instance of the worker pool config
func workerConfigLoader() {
	viper.SetDefault("WORKER_POOL_SIZE", 5)
	viper.SetDefault("WORKER_BUFFER_CHANNEL_SIZE", 100)
	viper.SetDefault("WORKER_BUFFER_FLUSH_TIMEOUT_MS", 5000)
	viper.SetDefault("WORKER_KAFKA_DELIVERY_CHANNEL_SIZE", 10)

	Worker = worker{
		WorkersPoolSize:     util.MustGetInt("WORKER_POOL_SIZE"),
		ChannelSize:         util.MustGetInt("WORKER_BUFFER_CHANNEL_SIZE"),
		DeliveryChannelSize: util.MustGetInt("WORKER_KAFKA_DELIVERY_CHANNEL_SIZE"),
		WorkerFlushTimeout:  util.MustGetDuration("WORKER_BUFFER_FLUSH_TIMEOUT_MS", time.Millisecond),
	}
}
