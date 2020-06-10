package publisher_test

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"raccoon/config"
	"raccoon/logger"
	"raccoon/publisher"
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
		kafkaproducer := &publisher.MockKafkaProducer{}
		kafkaproducer.On("Produce", mock.Anything, mock.Anything).Return(nil)
		publisher.NewProducer(&kafka.Producer{}, config.KafkaConfig{})
		err := kafkaproducer.Produce(&kafkaMessage, nil)

		assert.NoError(t, err)
	})

	suite.Run("MessageFailedToProduce", func(t *testing.T) {
		kafkaproducer := &publisher.MockKafkaProducer{}
		kafkaproducer.On("Produce", mock.Anything, mock.Anything).Return(fmt.Errorf("Error while producing into kafka"))
		publisher.NewProducer(&kafka.Producer{}, config.KafkaConfig{})
		err := kafkaproducer.Produce(&kafkaMessage, nil)

		assert.Error(t, err)
	})
}
