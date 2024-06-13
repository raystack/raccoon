package publisher

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/raystack/raccoon/metrics"
	pb "github.com/raystack/raccoon/proto"
)

// PubSub publishes to a Google Cloud Platform PubSub Topic.
type PubSub struct {
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

func (p *PubSub) ProduceBulk(events []*pb.Event, connGroup string) error {

	ctx := context.Background()
	errors := make([]error, len(events))
	results := make([]*pubsub.PublishResult, len(events))

	for order, event := range events {
		topicId := strings.Replace(p.topicFormat, "%s", event.Type, 1)

		topic, err := p.topic(ctx, topicId)
		if err != nil {
			metrics.Increment(
				"pubsub_messages_delivered_total",
				map[string]string{
					"success":    "false",
					"conn_group": connGroup,
					"event_type": event.Type,
				},
			)
			metrics.Increment(
				"pubsub_unknown_topic_failure_total",
				map[string]string{
					"topic":      topicId,
					"conn_group": connGroup,
					"event_type": event.Type,
				},
			)
			errors[order] = err
			continue
		}

		results[order] = topic.Publish(ctx, &pubsub.Message{
			Data: event.EventBytes,
		})

		metrics.Increment(
			"pubsub_messages_delivered_total",
			map[string]string{
				"success":    "true",
				"conn_group": connGroup,
				"event_type": event.Type,
			},
		)
	}

	for order, result := range results {
		if result == nil {
			continue
		}
		_, err := result.Get(ctx)
		if err != nil {
			metrics.Increment(
				"pubsub_messages_delivered_total",
				map[string]string{
					"success":    "false",
					"conn_group": connGroup,
					"event_type": events[order].Type,
				},
			)
			metrics.Increment(
				"pubsub_messages_undelivered_total",
				map[string]string{
					"success":    "true",
					"conn_group": connGroup,
					"event_type": events[order].Type,
				},
			)
			errors[order] = err
			continue
		}
	}

	if allNil(errors) {
		return nil
	}

	return BulkError{errors}
}

func (p *PubSub) topic(ctx context.Context, id string) (*pubsub.Topic, error) {

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
			return nil, fmt.Errorf(
				"topic %q doesn't exist in %q project", topic, p.client.Project(),
			)
		}

		cfg := &pubsub.TopicConfig{}
		if p.topicRetentionDuration > 0 {
			cfg.RetentionDuration = p.topicRetentionDuration
		}

		topic, err = p.client.CreateTopicWithConfig(ctx, id, cfg)
		if err != nil {
			return nil, fmt.Errorf("error creating topic %q: %w", id, err)
		}
		topic.PublishSettings = p.publishSettings
	}

	p.topics[id] = topic
	return topic, nil
}

func (p *PubSub) Close() error {
	p.topicLock.Lock()
	defer p.topicLock.Unlock()
	for _, topic := range p.topics {
		topic.Stop()
	}
	return p.client.Close()
}

func (p *PubSub) Name() string {
	return "pubsub"
}

type PubSubOpt func(*PubSub)

func WithPubSubTopicAutocreate(autocreate bool) PubSubOpt {
	return func(pub *PubSub) {
		pub.autoCreateTopic = autocreate
	}
}

func WithPubSubTopicRetention(duration time.Duration) PubSubOpt {
	return func(pub *PubSub) {
		pub.topicRetentionDuration = duration
	}
}

func WithPubSubDelayThreshold(duration time.Duration) PubSubOpt {
	return func(pub *PubSub) {
		pub.publishSettings.DelayThreshold = duration
	}
}

func WithPubSubCountThreshold(count int) PubSubOpt {
	return func(pub *PubSub) {
		pub.publishSettings.CountThreshold = count
	}
}

func WithPubSubByteThreshold(bytes int) PubSubOpt {
	return func(pub *PubSub) {
		pub.publishSettings.ByteThreshold = bytes
	}
}

func WithPubSubTimeout(timeout time.Duration) PubSubOpt {
	return func(pub *PubSub) {
		pub.publishSettings.Timeout = timeout
	}
}

func WithPubSubTopicFormat(format string) PubSubOpt {
	return func(pub *PubSub) {
		pub.topicFormat = format
	}
}

// NewPubSub creates a new PubSub publisher
func NewPubSub(client *pubsub.Client, opts ...PubSubOpt) (*PubSub, error) {

	p := &PubSub{
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
