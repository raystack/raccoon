package publisher_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	raccoonv1 "github.com/raystack/raccoon/proto"
	"github.com/raystack/raccoon/publisher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	envPubsubEmulator = "PUBSUB_EMULATOR_HOST"
	testingProject    = "test-project"
)

func TestPubSubPublisher(t *testing.T) {
	host := os.Getenv(envPubsubEmulator)
	if strings.TrimSpace(host) == "" {
		t.Logf(
			"skipping pubsub tests, because %s env variable is not set",
			envPubsubEmulator,
		)
		return
	}

	testEvent := &raccoonv1.Event{
		EventBytes: []byte("EVENT"),
		Type:       "click",
	}

	t.Run("should produce message successfully", func(t *testing.T) {
		client, err := pubsub.NewClient(context.Background(), testingProject)
		assert.NoError(t, err, "error creating pubsub client")

		pub, err := publisher.NewPubSub(
			client,
			publisher.WithPubSubTopicAutocreate(true),
			publisher.WithPubSubTopicRetentionDuration(10*time.Minute),
		)
		require.NoError(t, err, "unexpected error creating publisher")

		err = pub.ProduceBulk([]*raccoonv1.Event{testEvent}, "group")
		require.NoError(t, err, "error producing events")

		err = pub.Close()
		require.NoError(t, err, "error closing publisher")

		// publisher.Close() closes the client, so we create a new one
		client, err = pubsub.NewClient(context.Background(), testingProject)
		require.NoError(t, err)

		sub, err := client.CreateSubscription(
			context.Background(),
			"test-consumer",
			pubsub.SubscriptionConfig{
				Topic: client.Topic(testEvent.Type),
			},
		)
		require.NoError(t, err, "error creating subscription")

		ctx, cancel := context.WithCancel(context.Background())
		err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
			assert.Equal(t, testEvent.EventBytes, m.Data)
			m.Ack()
			cancel()
		})
		sub.Delete(context.Background())
		require.NoError(t, err, "error deleting subscription")

		err = client.Topic(testEvent.Type).Delete(context.Background())
		require.NoError(t, err, "error deleting topic")
	})

	t.Run("should return an error if topic doesn't exist and topic autocreate is set to false", func(t *testing.T) {

		client, err := pubsub.NewClient(context.Background(), testingProject)
		assert.NoError(t, err, "error creating pubsub client")

		pub, err := publisher.NewPubSub(client)
		require.NoError(t, err, "unexpected error creating publisher")

		err = pub.ProduceBulk([]*raccoonv1.Event{testEvent}, "group")
		require.Error(t, err)

		err = pub.Close()
		require.NoError(t, err, "error closing publisher")
	})
}
