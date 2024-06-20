package kinesis

import (
	"cmp"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	pb "github.com/raystack/raccoon/proto"
	"github.com/raystack/raccoon/publisher"
)

type Publisher struct {
	client              *kinesis.Client
	streamPattern       string
	streamAutocreate    bool
	streamProbeInterval time.Duration
	defaultShardCount   int32
}

func (p *Publisher) ProduceBulk(events []*pb.Event, connGroup string) error {
	errors := make([]error, len(events))
	for order, event := range events {
		streamName := strings.Replace(p.streamPattern, "%s", event.Type, 1)
		_, err := p.stream(streamName)
		if err != nil {
			errors[order] = err
			continue
		}

	}
	if cmp.Or(errors...) != nil {
		return &publisher.BulkError{Errors: errors}
	}
	return nil
}

func (p *Publisher) stream(name string) (string, error) {
	stream, err := p.client.DescribeStreamSummary(
		context.Background(),
		&kinesis.DescribeStreamSummaryInput{
			StreamName: aws.String(name),
		},
	)

	if err != nil {
		var errNotFound *types.ResourceNotFoundException
		if !errors.As(err, &errNotFound) {
			return "", err
		}

		if !p.streamAutocreate {
			return "", errNotFound
		}

		// TODO: handle ResourceLimitExceeded exception
		_, err := p.client.CreateStream(
			context.Background(),
			&kinesis.CreateStreamInput{
				StreamName: aws.String(name),
				ShardCount: aws.Int32(p.defaultShardCount),
				StreamModeDetails: &types.StreamModeDetails{
					StreamMode: types.StreamModeOnDemand, // TODO: make this configurable
				},
			},
		)
		if err != nil {
			return "", err
		}
		stream, err = p.client.DescribeStreamSummary(
			context.Background(),
			&kinesis.DescribeStreamSummaryInput{
				StreamName: aws.String(name),
			},
		)
		if err != nil {
			return "", err
		}
	}

	for stream.StreamDescriptionSummary.StreamStatus != types.StreamStatusActive {
		time.Sleep(p.streamProbeInterval)
		stream, err = p.client.DescribeStreamSummary(
			context.Background(),
			&kinesis.DescribeStreamSummaryInput{
				StreamName: aws.String(name),
			},
		)
		if err != nil {
			return "", err
		}
	}

	return *stream.StreamDescriptionSummary.StreamARN, nil
}

type Opt func(*Publisher)

func WithStreamAutocreate(autocreate bool) Opt {
	return func(p *Publisher) {
		p.streamAutocreate = autocreate
	}
}

func New(client *kinesis.Client, opts ...Opt) *Publisher {
	p := &Publisher{
		client:              client,
		streamPattern:       "%s",
		defaultShardCount:   1,
		streamProbeInterval: time.Second,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}
