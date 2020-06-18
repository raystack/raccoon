package publisher

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"raccoon/config"
	"raccoon/logger"
	"testing"
)

func init() {
	log, _ := test.NewNullLogger()
	logger.Set(log)
}

func TestProducer_Produce(suite *testing.T) {
	suite.Parallel()
	topic := "test_topic"
	kafkaMessage := kafka.Message{
		Key:   []byte("some_key"),
		Value: []byte("some_data"),
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
	}
	suite.Run("MessageSuccessfulProduce", func(t *testing.T) {
		kafkaproducer := &mockClient{}
		kafkaproducer.On("Produce", mock.Anything, mock.Anything).Return(nil)
		err := kafkaproducer.Produce(&kafkaMessage, nil)
		NewKafkaFromClient(kafkaproducer, config.KafkaConfig{})
		assert.NoError(t, err)
	})

	suite.Run("MessageFailedToProduce", func(t *testing.T) {
		kafkaproducer := &mockClient{}
		kafkaproducer.On("Produce", mock.Anything, mock.Anything).Return(fmt.Errorf("Error while producing into kafka"))
		err := kafkaproducer.Produce(&kafkaMessage, nil)
		NewKafkaFromClient(kafkaproducer, config.KafkaConfig{})
		assert.Error(t, err)
	})

	suite.Run("Should flush before closing the client", func(t *testing.T) {
		client := &mockClient{}
		client.On("Flush", 10).Return(0)
		client.On("Close").Return()
		kp := NewKafkaFromClient(client, config.KafkaConfig{
			FlushInterval: 10,
		})
		kp.Close()
		client.AssertExpectations(t)
	})
}
