package publisher

import (
	"fmt"
	"os"
	"testing"

	"github.com/odpf/raccoon/collection"
	"github.com/odpf/raccoon/logger"
	pb "github.com/odpf/raccoon/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type void struct{}

func (v void) Write(_ []byte) (int, error) {
	return 0, nil
}
func TestMain(t *testing.M) {
	logger.SetOutput(void{})
	os.Exit(t.Run())
}

func TestProducer_Close(suite *testing.T) {
	suite.Run("Should flush before closing the client", func(t *testing.T) {
		client := &mockClient{}
		client.On("Flush", 10).Return(0)
		client.On("Close").Return()
		kp := NewKafkaFromClient(client, 10, "%s")
		kp.Close()
		client.AssertExpectations(t)
	})
}

func TestKafka_ProduceBulk(suite *testing.T) {
	suite.Parallel()
	topic := "test_topic"
	suite.Run("AllMessagesSuccessfulProduce", func(t *testing.T) {
		t.Run("Should return nil when all message succesfully published", func(t *testing.T) {
			client := &mockClient{}
			client.On("Produce", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				go func() {
					args.Get(1).(chan kafka.Event) <- &kafka.Message{
						TopicPartition: kafka.TopicPartition{
							Topic:     args.Get(0).(*kafka.Message).TopicPartition.Topic,
							Partition: 0,
							Offset:    0,
							Error:     nil,
						},
					}
				}()
			})
			kp := NewKafkaFromClient(client, 10, "%s")

			err := kp.ProduceBulk(collection.CollectRequest{SendEventRequest: &pb.SendEventRequest{Events: []*pb.Event{{EventBytes: []byte{}, Type: topic}, {EventBytes: []byte{}, Type: topic}}}}, make(chan kafka.Event, 2))
			assert.NoError(t, err)
		})
	})

	suite.Run("PartialSuccessfulProduce", func(t *testing.T) {
		t.Run("Should process non producer error messages", func(t *testing.T) {
			client := &mockClient{}
			client.On("Produce", mock.Anything, mock.Anything).Return(fmt.Errorf("buffer full")).Once()
			client.On("Produce", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				go func() {
					args.Get(1).(chan kafka.Event) <- &kafka.Message{
						TopicPartition: kafka.TopicPartition{
							Topic:     args.Get(0).(*kafka.Message).TopicPartition.Topic,
							Partition: 0,
							Offset:    0,
							Error:     nil,
						},
					}
				}()
			}).Once()
			client.On("Produce", mock.Anything, mock.Anything).Return(fmt.Errorf("buffer full")).Once()
			kp := NewKafkaFromClient(client, 10, "%s")

			err := kp.ProduceBulk(collection.CollectRequest{SendEventRequest: &pb.SendEventRequest{Events: []*pb.Event{{EventBytes: []byte{}, Type: topic}, {EventBytes: []byte{}, Type: topic}, {EventBytes: []byte{}, Type: topic}}}}, make(chan kafka.Event, 2))
			assert.Len(t, err.(BulkError).Errors, 3)
			assert.Error(t, err.(BulkError).Errors[0])
			assert.Empty(t, err.(BulkError).Errors[1])
			assert.Error(t, err.(BulkError).Errors[2])
		})

		t.Run("Should return topic name when unknown topic is returned", func(t *testing.T) {
			client := &mockClient{}
			client.On("Produce", mock.Anything, mock.Anything).Return(fmt.Errorf("Local: Unknown topic")).Once()
			kp := NewKafkaFromClient(client, 10, "%s")

			err := kp.ProduceBulk(collection.CollectRequest{SendEventRequest: &pb.SendEventRequest{Events: []*pb.Event{{EventBytes: []byte{}, Type: topic}}}}, make(chan kafka.Event, 2))
			assert.EqualError(t, err.(BulkError).Errors[0], "Local: Unknown topic "+topic)
		})
	})

	suite.Run("MessageFailedToProduce", func(t *testing.T) {
		t.Run("Should fill all errors when all messages fail", func(t *testing.T) {
			client := &mockClient{}
			client.On("Produce", mock.Anything, mock.Anything).Return(fmt.Errorf("buffer full")).Once()
			client.On("Produce", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				go func() {
					args.Get(1).(chan kafka.Event) <- &kafka.Message{
						TopicPartition: kafka.TopicPartition{
							Topic:     args.Get(0).(*kafka.Message).TopicPartition.Topic,
							Partition: 0,
							Offset:    0,
							Error:     fmt.Errorf("timeout"),
						},
						Opaque: 1,
					}
				}()
			}).Once()
			kp := NewKafkaFromClient(client, 10, "%s")

			err := kp.ProduceBulk(collection.CollectRequest{SendEventRequest: &pb.SendEventRequest{Events: []*pb.Event{{EventBytes: []byte{}, Type: topic}, {EventBytes: []byte{}, Type: topic}}}}, make(chan kafka.Event, 2))
			assert.NotEmpty(t, err)
			assert.Len(t, err.(BulkError).Errors, 2)
			assert.Equal(t, "buffer full", err.(BulkError).Errors[0].Error())
			assert.Equal(t, "timeout", err.(BulkError).Errors[1].Error())
		})
	})
}
