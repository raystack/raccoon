package publisher

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	pb "github.com/raystack/raccoon/proto"
)

// PubSub publishes to a Google Cloud Platform PubSub Topic.
type PubSub struct {
	client      *pubsub.Client
	topicFormat string
}

func (p *PubSub) ProduceBulk(events []*pb.Event, connGroup string) error {

	// TODO(turtledev): instrument metrics

	ctx := context.Background()
	errors := make([]error, len(events))
	results := make([]*pubsub.PublishResult, len(events))

	// TODO(turtledev): topic cache can be shared across multiple ProduceBulk
	// invocations. But doing so introduces uncertainity with delivery guarantees.
	topics := make(map[string]*pubsub.Topic)

	for order, event := range events {
		topicId := fmt.Sprintf(p.topicFormat, event.Type)

		topic, exists := topics[topicId]
		if !exists {
			topic = p.client.Topic(topicId)
			valid, err := topic.Exists(ctx)
			if err != nil {
				return fmt.Errorf("error verifying existence of Topic %q: %w", topicId, err)
			}

			if !valid {
				topic, err = p.client.CreateTopic(ctx, topicId)
				// TODO(turtledev): guard against duplicate topic error
				if err != nil {
					return fmt.Errorf("error creating topic: %w", err)
				}
			}
			topics[topicId] = topic
		}

		results[order] = topic.Publish(ctx, &pubsub.Message{
			Data: event.EventBytes,
		})
	}

	for _, topic := range topics {
		topic.Stop()
	}

	for order, result := range results {
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

func (p *PubSub) Close() error {
	return p.client.Close()
}

func NewPubSub(projectId string, topicFormat string) (*PubSub, error) {
	// uses application default credentials
	// https://cloud.google.com/docs/authentication/application-default-credentials
	c, err := pubsub.NewClient(context.Background(), projectId)
	if err != nil {
		return nil, fmt.Errorf("NewPubSub: error creating client: %v", err)
	}

	return &PubSub{
		client:      c,
		topicFormat: topicFormat,
	}, nil
}
