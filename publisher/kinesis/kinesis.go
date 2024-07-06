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
	"github.com/raystack/raccoon/metrics"
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
				metrics.Increment(
					"kinesis_messages_delivered_total",
					map[string]string{
						"success":    "false",
						"conn_group": connGroup,
						"event_type": event.Type,
					},
				)
				if p.isErrNotFound(err) {
					metrics.Increment(
						"kinesis_unknown_stream_failure_total",
						map[string]string{
							"stream":     streamName,
							"conn_group": connGroup,
							"event_type": event.Type,
						},
					)
				}
				errors[order] = err
				continue
			}
		}

		metrics.Increment(
			"kinesis_messages_delivered_total",
			map[string]string{
				"success":    "true",
				"conn_group": connGroup,
				"event_type": event.Type,
			},
		)

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
			metrics.Increment(
				"kinesis_messages_delivered_total",
				map[string]string{
					"success":    "false",
					"conn_group": connGroup,
					"event_type": event.Type,
				},
			)
			metrics.Increment(
				"kinesis_messages_undelivered_total",
				map[string]string{
					"success":    "true",
					"conn_group": connGroup,
					"event_type": event.Type,
				},
			)
			if p.isErrNotFound(err) {
				metrics.Increment(
					"kinesis_unknown_stream_failure_total",
					map[string]string{
						"stream":     streamName,
						"conn_group": connGroup,
						"event_type": event.Type,
					},
				)
			}
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
		if !p.isErrNotFound(err) {
			return err
		}

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

func (*Publisher) isErrNotFound(e error) bool {
	var (
		errNotFound   *types.ResourceNotFoundException
		isErrNotFound = errors.As(e, &errNotFound)
	)
	return isErrNotFound
}

func (*Publisher) Name() string { return "kinesis" }
func (*Publisher) Close() error { return nil }

type Opt func(*Publisher) error

func WithStreamAutocreate(autocreate bool) Opt {
	return func(p *Publisher) error {
		p.streamAutocreate = autocreate
		return nil
	}
}

func WithStreamMode(mode types.StreamMode) Opt {

	validModesList := types.StreamMode("").Values()
	validModes := map[types.StreamMode]bool{}
	for _, m := range validModesList {
		validModes[types.StreamMode(m)] = true
	}

	return func(p *Publisher) error {
		valid := validModes[mode]
		if !valid {
			return fmt.Errorf(
				"unknown stream mode: %q (valid values: %v)",
				mode,
				validModesList,
			)
		}
		p.streamMode = mode
		return nil
	}
}

func WithShards(n uint32) Opt {
	return func(p *Publisher) error {
		p.defaultShardCount = int32(n)
		return nil
	}
}

func WithStreamPattern(pattern string) Opt {
	return func(p *Publisher) error {
		p.streamPattern = pattern
		return nil
	}
}

func WithPublishTimeout(timeout time.Duration) Opt {
	return func(p *Publisher) error {
		p.publishTimeout = timeout
		return nil
	}
}

func WithStreamProbleInterval(interval time.Duration) Opt {
	return func(p *Publisher) error {
		p.streamProbeInterval = interval
		return nil
	}
}

func New(client *kinesis.Client, opts ...Opt) (*Publisher, error) {
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
		err := opt(p)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}
