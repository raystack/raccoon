package kinesis

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	pb "github.com/raystack/raccoon/proto"
	"github.com/raystack/raccoon/publisher"
)

var globalCtx = context.Background()

type Publisher struct {
	client *kinesis.Client

	streamLock          sync.RWMutex
	streams             map[string]bool
	streamPattern       string
	streamAutocreate    bool
	streamProbeInterval time.Duration
	streamMode          types.StreamMode
	defaultShardCount   int32
	publishTimeout      time.Duration
}

func (p *Publisher) ProduceBulk(events []*pb.Event, connGroup string) error {

	ctx, cancel := context.WithTimeout(globalCtx, p.publishTimeout)
	defer cancel()
	errors := make([]error, len(events))

	for order, event := range events {
		streamName := strings.Replace(p.streamPattern, "%s", event.Type, 1)

		// only check for streams existence if publisher is configured
		// to create streams. If target stream doesn't exist, then
		// PutRecord will return an error anyway.
		if p.streamAutocreate {
			err := p.ensureStream(ctx, streamName)
			if err != nil {
				errors[order] = err
				continue
			}
		}

		partitionKey := fmt.Sprintf("%d", rand.Int31())
		_, err := p.client.PutRecord(
			ctx,
			&kinesis.PutRecordInput{
				Data:         event.EventBytes,
				PartitionKey: aws.String(partitionKey),
				StreamName:   aws.String(streamName),
			},
		)
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

func (p *Publisher) ensureStream(ctx context.Context, name string) error {

	p.streamLock.RLock()
	exists := p.streams[name]
	p.streamLock.RUnlock()

	if exists {
		return nil
	}
	p.streamLock.Lock()
	defer p.streamLock.Unlock()

	stream, err := p.client.DescribeStreamSummary(
		ctx,
		&kinesis.DescribeStreamSummaryInput{
			StreamName: aws.String(name),
		},
	)

	if err != nil {
		var errNotFound *types.ResourceNotFoundException
		if !errors.As(err, &errNotFound) {
			return err
		}

		if !p.streamAutocreate {
			return errNotFound
		}

		// TODO: handle ResourceLimitExceeded exception
		_, err := p.client.CreateStream(
			ctx,
			&kinesis.CreateStreamInput{
				StreamName: aws.String(name),
				ShardCount: aws.Int32(p.defaultShardCount),
				StreamModeDetails: &types.StreamModeDetails{
					StreamMode: p.streamMode,
				},
			},
		)
		if err != nil {
			return err
		}
		stream, err = p.client.DescribeStreamSummary(
			ctx,
			&kinesis.DescribeStreamSummaryInput{
				StreamName: aws.String(name),
			},
		)
		if err != nil {
			return err
		}
	}

	for stream.StreamDescriptionSummary.StreamStatus != types.StreamStatusActive {
		time.Sleep(p.streamProbeInterval)
		stream, err = p.client.DescribeStreamSummary(
			ctx,
			&kinesis.DescribeStreamSummaryInput{
				StreamName: aws.String(name),
			},
		)
		if err != nil {
			return err
		}
	}

	p.streams[name] = true
	return nil
}

func (*Publisher) Name() string { return "kinesis" }
func (*Publisher) Close() error { return nil }

type Opt func(*Publisher)

func WithStreamAutocreate(autocreate bool) Opt {
	return func(p *Publisher) {
		p.streamAutocreate = autocreate
	}
}

func WithStreamMode(mode types.StreamMode) Opt {
	return func(p *Publisher) {
		p.streamMode = mode
	}
}

func WithShards(n uint32) Opt {
	return func(p *Publisher) {
		p.defaultShardCount = int32(n)
	}
}

func WithStreamPattern(pattern string) Opt {
	return func(p *Publisher) {
		p.streamPattern = pattern
	}
}

func WithPublishTimeout(timeout time.Duration) Opt {
	return func(p *Publisher) {
		p.publishTimeout = timeout
	}
}

func New(client *kinesis.Client, opts ...Opt) *Publisher {
	p := &Publisher{
		client:              client,
		streamPattern:       "%s",
		defaultShardCount:   1,
		streamProbeInterval: time.Second,
		streamMode:          types.StreamModeOnDemand,
		streams:             make(map[string]bool),
		publishTimeout:      time.Minute,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}
