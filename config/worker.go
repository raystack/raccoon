package config

// WorkerConfig contains configs for kafka publisher worker pool
type WorkerConfig struct {
	workersPoolSize int
	channelSize     int
	outputTopic     string
}

// WorkerPoolSize number of worker to push to kafka initiated at the start of Raccoon
func (bc WorkerConfig) WorkersPoolSize() int {
	return bc.workersPoolSize
}

// ChannelSize channel size to buffer events before processed by worker
func (bc WorkerConfig) ChannelSize() int {
	return bc.channelSize
}

func (bc WorkerConfig) OutputTopic() string {
	return bc.outputTopic
}

func BufferConfigLoader() WorkerConfig {
	kc := WorkerConfig{
		workersPoolSize: mustGetInt("WORKER_POOL_SIZE"),
		channelSize:     mustGetInt("BUFFER_CHANNEL_SIZE"),
		outputTopic:     mustGetString("OUTPUT_TOPIC"),
	}
	return kc
}
