package kinesis_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	kinesis_sdk "github.com/aws/aws-sdk-go-v2/service/kinesis"
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

func TestKinesisProducer(t *testing.T) {

	localstackHost := os.Getenv(envLocalstackHost)
	if strings.TrimSpace(localstackHost) == "" {
		t.Errorf("cannot run tests because %s env variable is not set", envLocalstackHost)
		return
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	require.NoError(t, err, "error loading aws config")

	_ = kinesis_sdk.NewFromConfig(cfg, withLocalStack(localstackHost))
}
