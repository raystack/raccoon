package rest

import (
	"net/http"
	"time"

	"github.com/gojek/heimdall/v7"
	raccoon "github.com/goto/raccoon/clients/go"
	"github.com/goto/raccoon/clients/go/log"
	"github.com/goto/raccoon/clients/go/serializer"
	"github.com/goto/raccoon/clients/go/wire"
)

// Rest is the http implementation
type Rest struct {
	raccoon.Client
	url        string
	serialize  serializer.SerializerFunc
	wire       wire.WireMarshaler
	httpclient heimdall.Client
	headers    http.Header
	retryWait  time.Duration
	retryMax   uint
	logger     log.Logger
}

// Option represents the rest client options.
type Option func(*Rest)

// WithUrl sets the service address
func WithUrl(url string) Option {
	return func(rc *Rest) {
		rc.url = url
	}
}

// WithSerializer sets the serializer for the raccoon message.
func WithSerializer(s serializer.SerializerFunc) Option {
	return func(rc *Rest) {
		rc.serialize = s
	}
}

// WithHeader sets the http header for the request.
func WithHeader(key, val string) Option {
	return func(rc *Rest) {
		rc.headers.Add(key, val)
	}
}

// WithRetry retries for the error upto max attempts with the given delay between calls
func WithRetry(delay time.Duration, maxAttempts uint) Option {
	return func(rc *Rest) {
		rc.retryMax = maxAttempts
		rc.retryWait = delay
	}
}

// WithLogger sets the logger for the client.
func WithLogger(logger log.Logger) Option {
	return func(r *Rest) {
		log.SetLogger(logger)
		r.logger = log.Default()
	}
}
