package kinesis

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	pb "github.com/raystack/raccoon/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestKinesisProducer_UnitTest(t *testing.T) {
	events := []*pb.Event{
		{
			Type: "unknown",
		},
	}
	t.Run("should return an error if stream existence check fails", func(t *testing.T) {
		client := &mockKinesisClient{}

		client.On(
			"DescribeStreamSummary",
			mock.Anything,
			&kinesis.DescribeStreamSummaryInput{
				StreamName: aws.String("unknown"),
			},
			mock.Anything,
		).Return(
			&kinesis.DescribeStreamSummaryOutput{},
			fmt.Errorf("simulated error"),
		).Once()
		defer client.AssertExpectations(t)

		p, err := New(
			nil, // we will override it later
			WithStreamAutocreate(true),
		)
		if err != nil {
			t.Errorf("error constructing client: %v", err)
			return
		}
		p.client = client

		err = p.ProduceBulk(events, "")
		assert.Error(t, err, "error when sending message: simulated error")
	})
	t.Run("should return an error if stream creation exceeds resource limit", func(t *testing.T) {
		client := &mockKinesisClient{}

		client.On(
			"DescribeStreamSummary",
			mock.Anything,
			&kinesis.DescribeStreamSummaryInput{
				StreamName: aws.String("unknown"),
			},
			mock.Anything,
		).Return(
			&kinesis.DescribeStreamSummaryOutput{},
			&types.ResourceNotFoundException{},
		).Once()

		client.On("CreateStream", mock.Anything, mock.Anything, mock.Anything).
			Return(
				&kinesis.CreateStreamOutput{},
				&types.LimitExceededException{
					Message: aws.String("stream limit reached"),
				},
			).Once()
		defer client.AssertExpectations(t)

		p, err := New(
			nil, // we will override it later
			WithStreamAutocreate(true),
		)
		if err != nil {
			t.Errorf("error constructing client: %v", err)
			return
		}
		p.client = client

		err = p.ProduceBulk(events, "")
		assert.Error(t, err, "error when sending messages: LimitExceededException: stream limit reached")
	})
	t.Run("should return an error if rate limit is exceeded", func(t *testing.T) {

		client := &mockKinesisClient{}

		client.On("PutRecord", mock.Anything, mock.Anything, mock.Anything).
			Return(
				&kinesis.PutRecordOutput{},
				&types.ProvisionedThroughputExceededException{
					Message: aws.String("put limit exceeded"),
				},
			).Once()
		defer client.AssertExpectations(t)

		p, err := New(nil)
		if err != nil {
			t.Errorf("error constructing client: %v", err)
			return
		}
		p.client = client

		err = p.ProduceBulk(events, "")
		assert.Error(t, err, "error when sending messages: ProvisionedThroughputExceededException: put limit exceeded")
	})
}
