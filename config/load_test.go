package config

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	viper.Reset()
	viper.AutomaticEnv()
	os.Exit(m.Run())
}

func TestLogLevel(t *testing.T) {
	os.Setenv("LOG-LEVEL", "debug")
	logConfigLoader()
	assert.Equal(t, "debug", Log.Level)
}

func TestServerConfig(t *testing.T) {
	os.Setenv("SERVER-WEBSOCKET-PORT", "8080")
	os.Setenv("SERVER-WEBSOCKET-PING_INTERVAL", "1")
	os.Setenv("SERVER-WEBSOCKET-PONG_WAIT_INTERVAL", "1")
	os.Setenv("SERVER-WEBSOCKET-SERVER_SHUTDOWN_GRACE_PERIOD", "3")
	os.Setenv("SERVER-WEBSOKCET-USER_ID_HEADER", "x-user-id")
	serverConfigLoader()
	assert.Equal(t, "8080", Websocket.AppPort)
	assert.Equal(t, time.Duration(1)*time.Second, Websocket.PingInterval)
	assert.Equal(t, time.Duration(1)*time.Second, Websocket.PongWaitInterval)
	assert.Equal(t, time.Duration(3)*time.Second, Websocket.ServerShutDownGracePeriod)
}

func TestDynamicConfigLoad(t *testing.T) {
	os.Setenv("PUBLISHER-KAFKA-CLIENT-RANDOM", "anything")
	os.Setenv("PUBLISHER-KAFKA-CLIENT-BOOTSTRAP_SERVERS", "localhost:9092")
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(dynamicKafkaClientConfigLoad()))
	assert.Equal(t, "anything", viper.GetString("PUBLISHER-KAFKA-CLIENT-RANDOM"))
	assert.Equal(t, "localhost:9092", viper.GetString("PUBLISHER-KAFKA-CLIENT-BOOTSTRAP_SERVERS"))
}

func TestKafkaConfig_ToKafkaConfigMap(t *testing.T) {
	os.Setenv("PUBLISHER-KAFKA-FLUSH_INTERVAL", "1000")
	os.Setenv("PUBLISHER-KAFKA-CLIENT-BOOTSTRAP_SERVERS", "kafka:9092")
	os.Setenv("PUBLISHER-KAFKA-CLIENT-ACKS", "1")
	os.Setenv("PUBLISHER-KAFKA-CLIENT-QUEUE_BUFFERING_MAX_MESSAGES", "10000")
	os.Setenv("SOMETHING-PUBLISHER-KAFKA-CLIENT-SOMETHING", "anything")
	publisherConfigLoader()
	kafkaConfig := Kafka.ToKafkaConfigMap()
	bootstrapServer, _ := kafkaConfig.Get("bootstrap.servers", "")
	topic, _ := kafkaConfig.Get("topic", "")
	something, _ := kafkaConfig.Get("client.something", "")
	assert.Equal(t, "kafka:9092", bootstrapServer)
	assert.Equal(t, "", topic)
	assert.NotEqual(t, something, "anything")
	assert.Equal(t, 4, len(*kafkaConfig))
}

func TestWorkerConfig(t *testing.T) {
	os.Setenv("WORKER-POOL-SIZE", "2")
	os.Setenv("WORKER-BUFFER-CHANNEL_SIZE", "5")
	os.Setenv("WORKER-KAFKA-DELIVERY_CHANNEL_SIZE", "10")
	os.Setenv("WORKER-BUFFER-FLUSH_TIMEOUT", "100")
	workerConfigLoader()
	assert.Equal(t, time.Duration(100)*time.Second, Worker.WorkerFlushTimeout)
	assert.Equal(t, 10, Worker.DeliveryChannelSize)
	assert.Equal(t, 5, Worker.ChannelSize)
	assert.Equal(t, 2, Worker.WorkersPoolSize)
}
