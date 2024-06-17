package pubsub_test

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	pubsubsdk "cloud.google.com/go/pubsub"
	"github.com/raystack/raccoon/logger"
	raccoonv1 "github.com/raystack/raccoon/proto"
	"github.com/raystack/raccoon/publisher/pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	envPubsubEmulator = "PUBSUB_EMULATOR_HOST"
	testingProject    = "test-project"
)

var testEvent = &raccoonv1.Event{
	EventBytes: []byte("EVENT"),
	Type:       "click",
}

func TestPubSubPublisher(t *testing.T) {
	host := os.Getenv(envPubsubEmulator)
	if strings.TrimSpace(host) == "" {
		t.Logf(
			"skipping pubsub tests, because %s env variable is not set",
			envPubsubEmulator,
		)
		return
	}

	t.Run("should produce message successfully", func(t *testing.T) {
		client, err := pubsubsdk.NewClient(context.Background(), testingProject)
		assert.NoError(t, err, "error creating pubsub client")

		pub, err := pubsub.New(
			client,
			pubsub.WithTopicAutocreate(true),
			pubsub.WithTopicRetention(10*time.Minute),
		)
		require.NoError(t, err, "unexpected error creating publisher")

		err = pub.ProduceBulk([]*raccoonv1.Event{testEvent}, "group")
		require.NoError(t, err, "error producing events")

		err = pub.Close()
		require.NoError(t, err, "error closing publisher")

		// pub.Close() closes the client, so we create a new one
		client, err = pubsubsdk.NewClient(context.Background(), testingProject)
		require.NoError(t, err)

		sub, err := client.CreateSubscription(
			context.Background(),
			"test-consumer",
			pubsubsdk.SubscriptionConfig{
				Topic: client.Topic(testEvent.Type),
			},
		)
		require.NoError(t, err, "error creating subscription")

		ctx, cancel := context.WithCancel(context.Background())
		err = sub.Receive(ctx, func(ctx context.Context, m *pubsubsdk.Message) {
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

		client, err := pubsubsdk.NewClient(context.Background(), testingProject)
		assert.NoError(t, err, "error creating pubsub client")

		pub, err := pubsub.New(client)
		require.NoError(t, err, "unexpected error creating publisher")

		err = pub.ProduceBulk([]*raccoonv1.Event{testEvent}, "group")
		require.Error(t, err)

		err = pub.Close()
		require.NoError(t, err, "error closing publisher")
	})

	t.Run("should set retention for a topic correctly", func(t *testing.T) {

		client, err := pubsubsdk.NewClient(context.Background(), testingProject)
		assert.NoError(t, err, "error creating pubsub client")

		retention := time.Hour

		pub, err := pubsub.New(
			client,
			pubsub.WithTopicAutocreate(true),
			pubsub.WithTopicRetention(retention),
		)
		require.NoError(t, err, "unexpected error creating publisher")

		err = pub.ProduceBulk([]*raccoonv1.Event{testEvent}, "group")
		require.NoError(t, err, "error publishing events")

		cfg, err := client.Topic(testEvent.Type).Config(context.Background())
		require.NoError(t, err, "error obtaining topic config")
		require.Equal(t, cfg.RetentionDuration, retention)

		err = client.Topic(testEvent.Type).Delete(context.Background())
		require.NoError(t, err, "error deleting topic")

		err = pub.Close()
		require.NoError(t, err, "error closing publisher")
	})

	t.Run("should create the topic using topic format", func(t *testing.T) {

		client, err := pubsubsdk.NewClient(context.Background(), testingProject)
		assert.NoError(t, err, "error creating pubsub client")

		format := "pre-%s-post"
		pub, err := pubsub.New(
			client,
			pubsub.WithTopicAutocreate(true),
			pubsub.WithTopicFormat(format),
		)
		require.NoError(t, err, "unexpected error creating publisher")

		err = pub.ProduceBulk([]*raccoonv1.Event{testEvent}, "group")
		require.NoError(t, err, "error publishing events")

		topic := client.Topic("pre-click-post")
		exists, err := topic.Exists(context.Background())

		require.NoError(t, err, "error checking existence of topic")
		require.True(t, exists)

		err = topic.Delete(context.Background())
		require.NoError(t, err, "error deleting topic")

		err = pub.Close()
		require.NoError(t, err, "error closing publisher")
	})

	t.Run("static topic creation test", func(t *testing.T) {

		client, err := pubsubsdk.NewClient(context.Background(), testingProject)
		assert.NoError(t, err, "error creating pubsub client")

		format := "static-topic"
		pub, err := pubsub.New(
			client,
			pubsub.WithTopicAutocreate(true),
			pubsub.WithTopicFormat(format),
		)
		require.NoError(t, err, "unexpected error creating publisher")

		err = pub.ProduceBulk([]*raccoonv1.Event{testEvent}, "group")
		require.NoError(t, err, "error publishing events")

		topic := client.Topic(format)
		exists, err := topic.Exists(context.Background())

		require.NoError(t, err, "error checking existence of topic")
		require.True(t, exists)

		err = topic.Delete(context.Background())
		require.NoError(t, err, "error deleting topic")

		err = pub.Close()
		require.NoError(t, err, "error closing publisher")
	})
}

func TestMain(m *testing.M) {
	logger.SetOutput(io.Discard)
	os.Exit(m.Run())
}
