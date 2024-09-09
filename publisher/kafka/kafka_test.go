package kafka

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/raystack/raccoon/pkg/logger"
	pb "github.com/raystack/raccoon/proto"
	"github.com/raystack/raccoon/publisher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	group1 = "group-1"
)

func TestMain(t *testing.M) {
	logger.SetOutput(io.Discard)
	os.Exit(t.Run())
}

func TestProducer_Close(suite *testing.T) {
	suite.Run("Should flush before closing the client", func(t *testing.T) {
		client := &mockClient{}
		client.On("Flush", 10).Return(0)
		client.On("Close").Return()
		kp := NewFromClient(client, 10, "%s", 1)
		kp.Close()
		client.AssertExpectations(t)
	})
}

func TestKafka_ProduceBulk(suite *testing.T) {
	suite.Parallel()
	topic := "test_topic"
	suite.Run("AllMessagesSuccessfulProduce", func(t *testing.T) {
		t.Parallel()
		t.Run("Should return nil when all message successfully published", func(t *testing.T) {
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
						Opaque: 0,
					}
				}()
			})
			kp := NewFromClient(client, 10, "%s", 1)

			err := kp.ProduceBulk([]*pb.Event{{EventBytes: []byte{}, Type: topic}, {EventBytes: []byte{}, Type: topic}}, group1)
			assert.NoError(t, err)
		})
	})

	suite.Run("PartialSuccessfulProduce", func(t *testing.T) {
		t.Run("Should process non producer error messages", func(t *testing.T) {
			client := &mockClient{}
			client.On("Produce", mock.Anything, mock.Anything).Return(fmt.Errorf("buffer full")).Once()
			client.On("Produce", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				args.Get(1).(chan kafka.Event) <- &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic:     args.Get(0).(*kafka.Message).TopicPartition.Topic,
						Partition: 0,
						Offset:    0,
						Error:     nil,
					},
					Opaque: 1,
				}
			}).Once()
			client.On("Produce", mock.Anything, mock.Anything).Return(fmt.Errorf("buffer full")).Once()
			kp := NewFromClient(client, 10, "%s", 1)

			err := kp.ProduceBulk([]*pb.Event{{EventBytes: []byte{}, Type: topic}, {EventBytes: []byte{}, Type: topic}, {EventBytes: []byte{}, Type: topic}}, group1)
			assert.Len(t, err.(*publisher.BulkError).Errors, 3)
			assert.Error(t, err.(*publisher.BulkError).Errors[0])
			assert.Empty(t, err.(*publisher.BulkError).Errors[1])
			assert.Error(t, err.(*publisher.BulkError).Errors[2])
		})

		t.Run("Should return topic name when unknown topic is returned", func(t *testing.T) {
			client := &mockClient{}
			client.On("Produce", mock.Anything, mock.Anything).Return(fmt.Errorf("Local: Unknown topic")).Once()
			kp := NewFromClient(client, 10, "%s", 1)

			err := kp.ProduceBulk([]*pb.Event{{EventBytes: []byte{}, Type: topic}}, "group1")
			assert.EqualError(t, err.(*publisher.BulkError).Errors[0], "Local: Unknown topic "+topic)
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
			kp := NewFromClient(client, 10, "%s", 1)

			err := kp.ProduceBulk([]*pb.Event{{EventBytes: []byte{}, Type: topic}, {EventBytes: []byte{}, Type: topic}}, "group1")
			assert.NotEmpty(t, err)
			assert.Len(t, err.(*publisher.BulkError).Errors, 2)
			assert.Equal(t, "buffer full", err.(*publisher.BulkError).Errors[0].Error())
			assert.Equal(t, "timeout", err.(*publisher.BulkError).Errors[1].Error())
		})
	})
}
