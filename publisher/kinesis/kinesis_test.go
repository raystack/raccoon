package kinesis_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	kinesis_sdk "github.com/aws/aws-sdk-go-v2/service/kinesis"
	pb "github.com/raystack/raccoon/proto"
	"github.com/raystack/raccoon/publisher/kinesis"
	"github.com/stretchr/testify/require"
)

const (
	envLocalstackHost = "LOCALSTACK_HOST"
)

type localstackProvider struct{}

func (p *localstackProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     "test",
		SecretAccessKey: "test",
	}, nil
}

func withLocalStack(host string) func(o *kinesis_sdk.Options) {
	return func(o *kinesis_sdk.Options) {
		o.BaseEndpoint = aws.String(host)
		o.Credentials = &localstackProvider{}
	}
}

var (
	testEvent = &pb.Event{
		EventBytes: []byte("EVENT"),
		Type:       "click",
	}
)

func TestKinesisProducer(t *testing.T) {

	localstackHost := os.Getenv(envLocalstackHost)
	if strings.TrimSpace(localstackHost) == "" {
		t.Errorf("cannot run tests because %s env variable is not set", envLocalstackHost)
		return
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	require.NoError(t, err, "error loading aws config")

	client := kinesis_sdk.NewFromConfig(cfg, withLocalStack(localstackHost))

	t.Run("should return an error if stream doesn't exist", func(t *testing.T) {
		pub := kinesis.New(client)
		err := pub.ProduceBulk([]*pb.Event{testEvent}, "conn_group")
		require.Error(t, err)

	})
	t.Run("should create the stream if it doesn't exist and autocreate is set to true", func(t *testing.T) {
		pub := kinesis.New(client, kinesis.WithStreamAutocreate(true))
		err := pub.ProduceBulk([]*pb.Event{testEvent}, "conn_group")
		require.NoError(t, err)
		_, err = client.DeleteStream(context.Background(), &kinesis_sdk.DeleteStreamInput{
			StreamName: aws.String("click"),
		})
		require.NoError(t, err)
	})
}
