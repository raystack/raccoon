package config

var Worker worker

type buffer struct {
	// ChannelSize channel size to buffer events before processed by worker
	ChannelSize int `mapstructure:"channel_size" cmdx:"worker.buffer.channel.size" default:"100" desc:"Size of the buffer queue"`
	// FlushTimeoutMS specifies a timeout interval that the workers use to timeout
	// in case the workers could not complete the flush. This enables a non-blocking flush.
	FlushTimeoutMS int64 `mapstructure:"flush_timeout_ms" cmdx:"worker.buffer.flush.timeout.ms" default:"5000" desc:"Timeout for flushing leftover messages on shutdown"`
}

type worker struct {
	// WorkersPoolSize number of worker to push to kafka initiated at the start of Raccoon
	PoolSize int    `mapstructure:"pool_size" cmdx:"worker.pool.size" default:"5" desc:"No of workers that processes the events concurrently"`
	Buffer   buffer `mapstructure:"buffer"`

	//DeliveryChannelSize fetches the size of the delivery channel as configured
	DeliveryChannelSize int `mapstructure:"kafka_delivery_channel_size" cmdx:"worker.kafka.delivery.channel.size" default:"10"`
}
