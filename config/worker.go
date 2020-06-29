package config

import "github.com/spf13/viper"

var configLoaded bool
var wc WorkerConfig

// WorkerConfig contains configs for kafka publisher worker pool
type WorkerConfig struct {
	workersPoolSize     int
	channelSize         int
	deliveryChannelSize int
	workerFlushTimeout  int
}

// WorkersPoolSize number of worker to push to kafka initiated at the start of Raccoon
func (bc WorkerConfig) WorkersPoolSize() int {
	return bc.workersPoolSize
}

// ChannelSize channel size to buffer events before processed by worker
func (bc WorkerConfig) ChannelSize() int {
	return bc.channelSize
}

//DeliveryChannelSize fetches the size of the delivery channel as configured
func (bc WorkerConfig) DeliveryChannelSize() int {
	return bc.deliveryChannelSize
}

//WorkerFlushTimeout specifies a timeout interval that the workers use to timeout
//in case the workers could not complete the flush. This enables a non-blocking flush.
func (bc WorkerConfig) WorkerFlushTimeout() int {
	return bc.workerFlushTimeout
}

//WorkerConfigLoader constructs a singleton instance of the worker pool config
func WorkerConfigLoader() WorkerConfig {
	if !configLoaded {
		viper.SetDefault("WORKER_POOL_SIZE", 5)
		viper.SetDefault("BUFFER_CHANNEL_SIZE", 100)
		viper.SetDefault("DELIVERY_CHANNEL_SIZE", 10)
		viper.SetDefault("WORKER_FLUSH_TIMEOUT", 5)

		wc = WorkerConfig{
			workersPoolSize:     mustGetInt("WORKER_POOL_SIZE"),
			channelSize:         mustGetInt("BUFFER_CHANNEL_SIZE"),
			deliveryChannelSize: mustGetInt("DELIVERY_CHANNEL_SIZE"),
			workerFlushTimeout:  mustGetInt("WORKER_FLUSH_TIMEOUT"),
		}
	}
	return wc
}
