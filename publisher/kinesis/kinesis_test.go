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

	testCases := []struct {
		Desc        string
		Init        func(*mockKinesisClient)
		Opts        []Opt
		ExpectedErr string
	}{
		{
			Desc: "should return an error if stream existence check fails",
			Init: func(client *mockKinesisClient) {
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
			},
			Opts: []Opt{
				WithStreamAutocreate(true),
			},
			ExpectedErr: "error when sending message: simulated error",
		},
		{
			Desc: "should return an error if stream creation exceeds resource limit",
			Init: func(client *mockKinesisClient) {
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
			},
			Opts: []Opt{
				WithStreamAutocreate(true),
			},
			ExpectedErr: "error when sending messages: LimitExceededException: stream limit reached",
		},
		{
			Desc: "should return an error if rate limit is exceeded",
			Init: func(client *mockKinesisClient) {
				client.On("PutRecord", mock.Anything, mock.Anything, mock.Anything).
					Return(
						&kinesis.PutRecordOutput{},
						&types.ProvisionedThroughputExceededException{
							Message: aws.String("put limit exceeded"),
						},
					).Once()
			},
			ExpectedErr: "error when sending messages: ProvisionedThroughputExceededException: put limit exceeded",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Desc, func(t *testing.T) {
			client := &mockKinesisClient{}
			testCase.Init(client)
			defer client.AssertExpectations(t)

			p, err := New(client, testCase.Opts...)
			if err != nil {
				t.Errorf("error constructing client: %v", err)
				return
			}

			err = p.ProduceBulk(events, "")
			assert.Error(t, err, testCase.ExpectedErr)
		})

	}
}
