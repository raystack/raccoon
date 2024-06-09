package publisher

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	pb "github.com/raystack/raccoon/proto"
)

// PubSub publishes to a Google Cloud Platform PubSub Topic.
type PubSub struct {
	client *pubsub.Client
}

func (p *PubSub) ProduceBulk(events []*pb.Event, connGroup string) error {
	return nil
}

func (p *PubSub) Close() error {
	return p.client.Close()
}

func NewPubSub(projectId string) (*PubSub, error) {
	// uses applicated default credentials
	// https://cloud.google.com/docs/authentication/application-default-credentials
	c, err := pubsub.NewClient(context.Background(), projectId)
	if err != nil {
		return nil, fmt.Errorf("NewPubSub: error creating client: %v", err)
	}
	return &PubSub{c}, nil
}
