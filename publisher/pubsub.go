package publisher

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
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
	topicLock sync.RWMutex
	topics    map[string]*pubsub.Topic
}

func (p *PubSub) ProduceBulk(events []*pb.Event, connGroup string) error {

	// TODO(turtledev): instrument metrics

	ctx := context.Background()
	errors := make([]error, len(events))
	results := make([]*pubsub.PublishResult, len(events))

	for order, event := range events {
		topicId := fmt.Sprintf(p.topicFormat, event.Type)

		topic, err := p.topic(ctx, topicId)
		if err != nil {
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
		_, err := result.Get(ctx)
		if err != nil {
			errors[order] = err
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
		return nil, fmt.Errorf("error verifying existence of Topic %q: %w", id, err)
	}

	if !exists {
		topic, err = p.client.CreateTopic(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("error creating topic: %w", err)
		}
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

// NewPubSub creates a new PubSub publisher
// uses default application credentials
// https://cloud.google.com/docs/authentication/application-default-credentials
func NewPubSub(projectId string, topicFormat string) (*PubSub, error) {
	c, err := pubsub.NewClient(context.Background(), projectId)
	if err != nil {
		return nil, fmt.Errorf("NewPubSub: error creating client: %v", err)
	}

	return &PubSub{
		client:      c,
		topicFormat: topicFormat,
		topicLock:   sync.RWMutex{},
		topics:      make(map[string]*pubsub.Topic),
	}, nil
}
