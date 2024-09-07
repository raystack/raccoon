package kinesis

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	pb "github.com/raystack/raccoon/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestKinesisProducer_UnitTest(t *testing.T) {
	t.Run("should return an error if stream creation fails", func(t *testing.T) {
		events := []*pb.Event{
			{
				Type: "unknown",
			},
		}
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
		assert.NotNil(t, err)
	})
}
