package config

type worker struct {
	// WorkersPoolSize number of worker to push to kafka initiated at the start of Raccoon
	WorkersPoolSize int `mapstructure:"WORKER_POOL_SIZE" cmdx:"worker.pool.size" default:"5"`
	// ChannelSize channel size to buffer events before processed by worker
	ChannelSize int `mapstructure:"WORKER_BUFFER_CHANNEL_SIZE" cmdx:"worker.buffer.channel.size" default:"100"`
	//DeliveryChannelSize fetches the size of the delivery channel as configured
	DeliveryChannelSize int `mapstructure:"WORKER_KAFKA_DELIVERY_CHANNEL_SIZE" cmdx:"worker.kafka.delivery.channel.size" default:"10"`
	//WorkerFlushTimeoutMS specifies a timeout interval that the workers use to timeout
	//in case the workers could not complete the flush. This enables a non-blocking flush.
	WorkerFlushTimeoutMS int64 `mapstructure:"WORKER_BUFFER_FLUSH_TIMEOUT_MS" cmdx:"worker.buffer.flush.timeout.ms" default:"5000"`
}
