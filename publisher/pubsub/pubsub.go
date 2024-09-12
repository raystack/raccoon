package pubsub

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/googleapis/gax-go/v2/apierror"
	"github.com/raystack/raccoon/pkg/metrics"
	pb "github.com/raystack/raccoon/proto"
	"github.com/raystack/raccoon/publisher"
	"google.golang.org/grpc/codes"
)

var globalCtx = context.Background()

// Publisher publishes to a Google Cloud Platform PubSub Topic.
type Publisher struct {
	client      *pubsub.Client
	topicFormat string

	// TODO(turtledev): There is scope for optimising topic cache
	// Problems:
	//  * topic cache grows unbounded
	//  * every time a topic is added, it acquires a global lock.
	//    This causes readers of other topics to get blocked as well,
	//    which is not optimal for performance.
	topicLock              sync.RWMutex
	topics                 map[string]*pubsub.Topic
	autoCreateTopic        bool
	topicRetentionDuration time.Duration
	publishSettings        pubsub.PublishSettings
}

func (p *Publisher) ProduceBulk(events []*pb.Event, connGroup string) error {
	ctx, cancel := context.WithTimeout(globalCtx, p.publishSettings.Timeout)
	defer cancel()

	errors := make([]error, len(events))
	results := make([]*pubsub.PublishResult, len(events))

	for order, event := range events {
		topicName := p.topicNameFromEvent(event)

		topic, err := p.topic(ctx, topicName)
		if err != nil {
			reportTopicError(err, topicName, connGroup, event.Type)
			errors[order] = err
			continue
		}

		results[order] = topic.Publish(ctx, &pubsub.Message{
			Data: event.EventBytes,
		})
	}

	for order, result := range results {
		if result == nil {
			continue
		}
		var (
			event     = events[order]
			topicName = p.topicNameFromEvent(event)
		)
		_, err := result.Get(ctx)
		if err != nil {
			reportPublishError(err, topicName, connGroup, event.Type)
			errors[order] = err
			continue
		}

		metrics.Increment(
			"pubsub_messages_delivered_total",
			map[string]string{
				"topic":      topicName,
				"conn_group": connGroup,
				"event_type": events[order].Type,
			},
		)
	}

	if cmp.Or(errors...) != nil {
		return &publisher.BulkError{Errors: errors}
	}
	return nil
}

func (p *Publisher) topicNameFromEvent(event *pb.Event) string {
	return strings.Replace(p.topicFormat, "%s", event.Type, 1)
}

func (p *Publisher) topic(ctx context.Context, id string) (*pubsub.Topic, error) {
	p.topicLock.RLock()
	topic, exists := p.topics[id]
	p.topicLock.RUnlock()

	if exists {
		return topic, nil
	}

	p.topicLock.Lock()
	defer p.topicLock.Unlock()

	// double-checked locking in case another goroutine was blocked
	// on the same topic.
	// This ensures that we don't create duplicate topic instances.
	if p.topics[id] != nil {
		return p.topics[id], nil
	}

	topic = p.client.Topic(id)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("error verifying existence of topic %q: %w", id, err)
	}

	if !exists {
		if !p.autoCreateTopic {
			return nil, unknownTopicError{Topic: id, Project: p.client.Project()}
		}

		cfg := &pubsub.TopicConfig{}
		if p.topicRetentionDuration > 0 {
			cfg.RetentionDuration = p.topicRetentionDuration
		}

		topic, err = p.client.CreateTopicWithConfig(ctx, id, cfg)
		if err != nil {
			// in case a service replica created this topic before we could
			if isErrAlreadyExists(err) {
				topic = p.client.Topic(id)
			} else {
				return nil, fmt.Errorf("error creating topic %q: %w", id, err)
			}
		}
	}

	topic.PublishSettings = p.publishSettings
	p.topics[id] = topic
	return topic, nil
}

func (p *Publisher) Close() error {
	p.topicLock.Lock()
	defer p.topicLock.Unlock()
	for _, topic := range p.topics {
		topic.Stop()
	}
	return p.client.Close()
}

func (p *Publisher) Name() string {
	return "pubsub"
}

type Opt func(*Publisher)

func WithTopicAutocreate(autocreate bool) Opt {
	return func(pub *Publisher) {
		pub.autoCreateTopic = autocreate
	}
}

func WithTopicRetention(duration time.Duration) Opt {
	return func(pub *Publisher) {
		pub.topicRetentionDuration = duration
	}
}

func WithDelayThreshold(duration time.Duration) Opt {
	return func(pub *Publisher) {
		pub.publishSettings.DelayThreshold = duration
	}
}

func WithCountThreshold(count int) Opt {
	return func(pub *Publisher) {
		pub.publishSettings.CountThreshold = count
	}
}

func WithByteThreshold(bytes int) Opt {
	return func(pub *Publisher) {
		pub.publishSettings.ByteThreshold = bytes
	}
}

func WithTimeout(timeout time.Duration) Opt {
	return func(pub *Publisher) {
		pub.publishSettings.Timeout = timeout
	}
}

func WithTopicFormat(format string) Opt {
	return func(pub *Publisher) {
		pub.topicFormat = format
	}
}

// NewPubSub creates a new PubSub publisher
func New(client *pubsub.Client, opts ...Opt) (*Publisher, error) {
	p := &Publisher{
		client:          client,
		topicFormat:     "%s",
		topicLock:       sync.RWMutex{},
		topics:          make(map[string]*pubsub.Topic),
		publishSettings: pubsub.DefaultPublishSettings,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p, nil
}

func isErrUnknownTopic(e error) bool {
	return errors.As(e, &unknownTopicError{})
}

func isErrAlreadyExists(e error) bool {
	return hasAPIErrorCode(e, codes.AlreadyExists)
}

func isErrResourceExhausted(e error) bool {
	return hasAPIErrorCode(e, codes.ResourceExhausted)
}

func hasAPIErrorCode(e error, code codes.Code) bool {
	apiError, ok := e.(*apierror.APIError)
	if !ok {
		return false
	}
	return apiError.GRPCStatus().Code() == code
}

func reportTopicError(err error, topicName, connGroup, eventType string) {
	metrics.Increment(
		"pubsub_messages_undelivered_total",
		map[string]string{
			"topic":      topicName,
			"conn_group": connGroup,
			"event_type": eventType,
		},
	)
	switch {
	case isErrUnknownTopic(err):
		metrics.Increment(
			"pubsub_unknown_topic_failure_total",
			map[string]string{
				"topic":      topicName,
				"conn_group": connGroup,
				"event_type": eventType,
			},
		)
	case isErrResourceExhausted(err):
		metrics.Increment(
			"pubsub_topics_limit_exceeded_total",
			map[string]string{
				"topic":      topicName,
				"conn_group": connGroup,
				"event_type": eventType,
			},
		)
	}
}

func reportPublishError(err error, topicName, connGroup, eventType string) {
	metrics.Increment(
		"pubsub_messages_undelivered_total",
		map[string]string{
			"topic":      topicName,
			"conn_group": connGroup,
			"event_type": eventType,
		},
	)
	if isErrResourceExhausted(err) {
		metrics.Increment(
			"pubsub_topic_throughput_exceeded_total",
			map[string]string{
				"topic":      topicName,
				"conn_group": connGroup,
				"event_type": eventType,
			},
		)
	}
}
