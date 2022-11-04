package grpc

import (
	"time"

	"github.com/odpf/raccoon/clients/go/log"

	"google.golang.org/grpc"

	raccoon "github.com/odpf/raccoon/clients/go"
	"github.com/odpf/raccoon/clients/go/serializer"
)

// Grpc is the http implementation
type Grpc struct {
	raccoon.Client
	addr        string
	client      *grpc.ClientConn
	serialize   serializer.SerializerFunc
	headers     map[string]string
	dialOptions []grpc.DialOption
	retryMax    uint
	retryWait   time.Duration
	logger      log.Logger
}

// Option represents the grpc client options.
type Option func(*Grpc)

// WithAddr sets the service address
func WithAddr(addr string) Option {
	return func(gc *Grpc) {
		gc.addr = addr
	}
}

// WithSerializer sets the serializer for the raccoon message.
func WithSerializer(s serializer.SerializerFunc) Option {
	return func(gc *Grpc) {
		gc.serialize = s
	}
}

// WithHeader sets the grpc metadata for the request.
func WithHeader(key, val string) Option {
	return func(gc *Grpc) {
		gc.headers[key] = val
	}
}

// WithDialOptions sets the grpc dial options.
func WithDialOptions(opts ...grpc.DialOption) Option {
	return func(gc *Grpc) {
		gc.dialOptions = opts
	}
}

// WithRetry retries for the error upto max attempts with the given delay between calls
func WithRetry(delay time.Duration, maxAttempts uint) Option {
	return func(gc *Grpc) {
		gc.retryMax = maxAttempts
		gc.retryWait = delay
	}
}

// WithLogger sets the logger for the client.
func WithLogger(logger log.Logger) Option {
	return func(gc *Grpc) {
		log.SetLogger(logger)
		gc.logger = log.Default()
	}
}
