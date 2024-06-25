package kinesis_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	kinesis_sdk "github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	pb "github.com/raystack/raccoon/proto"
	"github.com/raystack/raccoon/publisher/kinesis"
	"github.com/stretchr/testify/require"
)

const (
	envLocalstackHost = "LOCALSTACK_HOST"
)

type localstackProvider struct{}

func (p *localstackProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     "test",
		SecretAccessKey: "test",
	}, nil
}

func withLocalStack(host string) func(o *kinesis_sdk.Options) {
	return func(o *kinesis_sdk.Options) {
		o.BaseEndpoint = aws.String(host)
		o.Credentials = &localstackProvider{}
	}
}

var (
	testEvent = &pb.Event{
		EventBytes: []byte("EVENT"),
		Type:       "click",
	}
)

func createStream(client *kinesis_sdk.Client, name string) (string, error) {
	_, err := client.CreateStream(
		context.Background(),
		&kinesis_sdk.CreateStreamInput{
			StreamName: aws.String(name),
			StreamModeDetails: &types.StreamModeDetails{
				StreamMode: types.StreamModeOnDemand,
			},
			ShardCount: aws.Int32(1),
		},
	)
	if err != nil {
		return "", err
	}
	retries := 5
	for range retries {
		stream, err := client.DescribeStreamSummary(
			context.Background(),
			&kinesis_sdk.DescribeStreamSummaryInput{
				StreamName: aws.String(name),
			},
		)
		if err != nil {
			return "", err
		}
		if stream.StreamDescriptionSummary.StreamStatus == types.StreamStatusActive {
			return *stream.StreamDescriptionSummary.StreamARN, nil
		}
		time.Sleep(time.Second / 2)
	}
	return "", fmt.Errorf("timed out waiting for stream to get ready")
}

func deleteStream(client *kinesis_sdk.Client, name string) error {
	_, err := client.DeleteStream(context.Background(), &kinesis_sdk.DeleteStreamInput{
		StreamName: aws.String(name),
	})
	if err != nil {
		return err
	}

	var errNotFound *types.ResourceNotFoundException
	for !errors.As(err, &errNotFound) {
		_, err = client.DescribeStreamSummary(
			context.Background(),
			&kinesis_sdk.DescribeStreamSummaryInput{
				StreamName: aws.String(name),
			},
		)
		time.Sleep(time.Second / 2)
	}

	return nil
}

func getStreamMode(client *kinesis_sdk.Client, name string) (types.StreamMode, error) {
	stream, err := client.DescribeStreamSummary(
		context.Background(),
		&kinesis_sdk.DescribeStreamSummaryInput{
			StreamName: aws.String(name),
		},
	)
	if err != nil {
		return "", err
	}
	return stream.StreamDescriptionSummary.StreamModeDetails.StreamMode, nil
}

func readStream(client *kinesis_sdk.Client, arn string) ([][]byte, error) {
	stream, err := client.DescribeStream(
		context.Background(),
		&kinesis_sdk.DescribeStreamInput{
			StreamARN: aws.String(arn),
		},
	)
	if err != nil {
		return nil, err
	}
	if len(stream.StreamDescription.Shards) == 0 {
		return nil, fmt.Errorf("stream %q has no shards", arn)
	}
	iter, err := client.GetShardIterator(
		context.Background(),
		&kinesis_sdk.GetShardIteratorInput{
			ShardId:           stream.StreamDescription.Shards[0].ShardId,
			StreamARN:         aws.String(arn),
			ShardIteratorType: types.ShardIteratorTypeTrimHorizon,
		},
	)
	if err != nil {
		return nil, err
	}
	res, err := client.GetRecords(
		context.Background(),
		&kinesis_sdk.GetRecordsInput{
			StreamARN:     aws.String(arn),
			ShardIterator: iter.ShardIterator,
		},
	)
	if err != nil {
		return nil, err
	}
	if len(res.Records) == 0 {
		return nil, fmt.Errorf("got empty response")
	}
	rv := [][]byte{}
	for _, record := range res.Records {
		rv = append(rv, record.Data)
	}
	return rv, nil
}

func TestKinesisProducer(t *testing.T) {

	localstackHost := os.Getenv(envLocalstackHost)
	if strings.TrimSpace(localstackHost) == "" {
		t.Errorf("cannot run tests because %s env variable is not set", envLocalstackHost)
		return
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	require.NoError(t, err, "error loading aws config")

	client := kinesis_sdk.NewFromConfig(cfg, withLocalStack(localstackHost))

	t.Run("should return an error if stream doesn't exist", func(t *testing.T) {
		pub, err := kinesis.New(client)
		require.NoError(t, err)
		err = pub.ProduceBulk([]*pb.Event{testEvent}, "conn_group")
		require.Error(t, err)

	})

	t.Run("should return an error if an invalid stream mode is specified", func(t *testing.T) {
		_, err := kinesis.New(
			client,
			kinesis.WithStreamMode("INVALID"),
		)
		require.Error(t, err)
	})

	t.Run("should publish message to kinesis", func(t *testing.T) {
		streamARN, err := createStream(client, testEvent.Type)
		require.NoError(t, err)
		defer deleteStream(client, testEvent.Type)

		pub, err := kinesis.New(client)
		require.NoError(t, err)
		pub.ProduceBulk([]*pb.Event{testEvent}, "conn_group")
		require.NoError(t, err)
		events, err := readStream(client, streamARN)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, events[0], testEvent.EventBytes)
	})
	t.Run("stream auto creation", func(t *testing.T) {
		t.Run("should create the stream if it doesn't exist and autocreate is set to true", func(t *testing.T) {
			pub, err := kinesis.New(client, kinesis.WithStreamAutocreate(true))
			require.NoError(t, err)

			err = pub.ProduceBulk([]*pb.Event{testEvent}, "conn_group")
			require.NoError(t, err)
			deleteStream(client, testEvent.Type)
		})
		t.Run("should create the stream with mode = ON_DEMAND (default)", func(t *testing.T) {

			pub, err := kinesis.New(client, kinesis.WithStreamAutocreate(true))
			require.NoError(t, err)
			err = pub.ProduceBulk([]*pb.Event{testEvent}, "conn_group")
			require.NoError(t, err)
			defer deleteStream(client, testEvent.Type)

			mode, err := getStreamMode(client, testEvent.Type)
			require.NoError(t, err)
			require.Equal(t, mode, types.StreamModeOnDemand)
		})
		t.Run("should create the stream with mode = PROVISIONED", func(t *testing.T) {
			pub, err := kinesis.New(
				client,
				kinesis.WithStreamAutocreate(true),
				kinesis.WithStreamMode(types.StreamModeProvisioned),
			)
			require.NoError(t, err)
			err = pub.ProduceBulk([]*pb.Event{testEvent}, "conn_group")
			require.NoError(t, err)
			defer deleteStream(client, testEvent.Type)

			mode, err := getStreamMode(client, testEvent.Type)
			require.NoError(t, err)
			require.Equal(t, mode, types.StreamModeProvisioned)
		})
		t.Run("should create stream with specified number of shards", func(t *testing.T) {
			shards := 5
			pub, err := kinesis.New(
				client,
				kinesis.WithStreamAutocreate(true),
				kinesis.WithShards(uint32(shards)),
			)
			require.NoError(t, err)

			err = pub.ProduceBulk([]*pb.Event{testEvent}, "conn_group")
			require.NoError(t, err)
			defer deleteStream(client, testEvent.Type)

			stream, err := client.DescribeStream(
				context.Background(),
				&kinesis_sdk.DescribeStreamInput{
					StreamName: aws.String(testEvent.Type),
				},
			)
			require.NoError(t, err)
			require.Equal(t, shards, len(stream.StreamDescription.Shards))
		})
	})

	t.Run("should publish message according to the stream pattern", func(t *testing.T) {
		streamPattern := "pre-%s-post"
		destinationStream := "pre-click-post"
		_, err := createStream(client, destinationStream)
		require.NoError(t, err)
		defer deleteStream(client, destinationStream)
		pub, err := kinesis.New(
			client,
			kinesis.WithStreamPattern(streamPattern),
		)
		require.NoError(t, err)
		err = pub.ProduceBulk([]*pb.Event{testEvent}, "conn_group")
		require.NoError(t, err)
	})
	t.Run("should publish messages to static stream names", func(t *testing.T) {
		destinationStream := "static"
		_, err := createStream(client, destinationStream)
		require.NoError(t, err)
		defer deleteStream(client, destinationStream)
		pub, err := kinesis.New(
			client,
			kinesis.WithStreamPattern(destinationStream),
		)
		require.NoError(t, err)
		err = pub.ProduceBulk([]*pb.Event{testEvent}, "conn_group")
		require.NoError(t, err)
	})
}
