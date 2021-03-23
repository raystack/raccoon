package config

import (
	"raccoon/config/util"
	"time"

	"github.com/spf13/viper"
)

// Worker contains configs for kafka publisher worker pool
var Worker worker

type worker struct {
	// WorkersPoolSize number of worker to push to kafka initiated at the start of Raccoon
	WorkersPoolSize int
	// ChannelSize channel size to buffer events before processed by worker
	ChannelSize int
	//DeliveryChannelSize fetches the size of the delivery channel as configured
	DeliveryChannelSize int
	//WorkerFlushTimeout specifies a timeout interval that the workers use to timeout
	//in case the workers could not complete the flush. This enables a non-blocking flush.
	WorkerFlushTimeout time.Duration
}

//workerConfigLoader constructs a singleton instance of the worker pool config
func workerConfigLoader() {
	viper.SetDefault("WORKER-POOL-SIZE", 5)
	viper.SetDefault("WORKER-BUFFER-CHANNEL_SIZE", 100)
	viper.SetDefault("WORKER-BUFFER-FLUSH_TIMEOUT", 5)
	viper.SetDefault("WORKER-KAFKA-DELIVERY_CHANNEL_SIZE", 10)

	Worker = worker{
		WorkersPoolSize:     util.MustGetInt("WORKER-POOL-SIZE"),
		ChannelSize:         util.MustGetInt("WORKER-BUFFER-CHANNEL_SIZE"),
		DeliveryChannelSize: util.MustGetInt("WORKER-KAFKA-DELIVERY_CHANNEL_SIZE"),
		WorkerFlushTimeout:  util.MustGetDuration("WORKER-BUFFER-FLUSH_TIMEOUT", time.Second),
	}
}
