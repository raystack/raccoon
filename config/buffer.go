package config

// BufferConfig contains configs for kafka publisher worker pool
type BufferConfig struct {
	workersPoolSize int
	channelSize     int
}

// WorkerPoolSize number of worker to push to kafka initiated at the start of Raccoon
func (bc BufferConfig) WorkersPoolSize() int {
	return bc.workersPoolSize
}

// ChannelSize channel size to buffer events before processed by worker
func (bc BufferConfig) ChannelSize() int {
	return bc.channelSize
}

func BufferConfigLoader() BufferConfig {
	kc := BufferConfig{
		workersPoolSize: mustGetInt("WORKER_POOL_SIZE"),
		channelSize:     mustGetInt("BUFFER_CHANNEL_SIZE"),
	}
	return kc
}
