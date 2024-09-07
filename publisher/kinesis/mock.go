package kinesis

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/stretchr/testify/mock"
)

type mockKinesisClient struct {
	mock.Mock
}

func (cli *mockKinesisClient) PutRecord(ctx context.Context, in *kinesis.PutRecordInput, opts ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error) {
	args := cli.Called(ctx, in, opts)
	return args.Get(0).(*kinesis.PutRecordOutput), args.Error(1)
}

func (cli *mockKinesisClient) DescribeStreamSummary(ctx context.Context, in *kinesis.DescribeStreamSummaryInput, opts ...func(*kinesis.Options)) (*kinesis.DescribeStreamSummaryOutput, error) {
	args := cli.Called(ctx, in, opts)
	return args.Get(0).(*kinesis.DescribeStreamSummaryOutput), args.Error(1)
}

func (cli *mockKinesisClient) CreateStream(ctx context.Context, in *kinesis.CreateStreamInput, opts ...func(*kinesis.Options)) (*kinesis.CreateStreamOutput, error) {
	args := cli.Called(ctx, in, opts)
	return args.Get(0).(*kinesis.CreateStreamOutput), args.Error(1)
}
