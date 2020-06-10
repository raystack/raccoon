package config

// BufferConfig contains configs for kafka publisher worker pool
type BufferConfig struct {
	poolNumbers int
	channelSize int
}

// PoolNumbers number of worker to push to kafka initiated at the start of Raccoon
func (bc BufferConfig) PoolNumbers() int {
	return bc.poolNumbers
}

// ChannelSize channel size to buffer events before processed by worker
func (bc BufferConfig) ChannelSize() int {
	return bc.channelSize
}

func BufferConfigLoader() BufferConfig {
	kc := BufferConfig{
		poolNumbers: mustGetInt("BUFFER_POOL_NUMBERS"),
		channelSize: mustGetInt("BUFFER_CHANNEL_SIZE"),
	}
	return kc
}
