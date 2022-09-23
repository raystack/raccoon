package rest

import (
	"net/http"

	"github.com/gojek/heimdall/v7"

	raccoon "github.com/odpf/raccoon/clients/go"
	"github.com/odpf/raccoon/clients/go/serializer"
	"github.com/odpf/raccoon/clients/go/wire"
)

// RestClient is the http implementation
type RestClient struct {
	raccoon.Client
	url        string
	serialize  serializer.SerializerFunc
	wire       wire.WireMarshaler
	httpclient heimdall.Client
	headers    http.Header
}

// RestOption represents the client options.
type RestOption func(*RestClient)

// WithUrl sets the service address
func WithUrl(url string) RestOption {
	return func(rc *RestClient) {
		rc.url = url
	}
}

// WithSerializer sets the serializer for the raccoon message.
func WithSerializer(s serializer.SerializerFunc) RestOption {
	return func(rc *RestClient) {
		rc.serialize = s
	}
}

// WithHeader sets the http header for the request.
func WithHeader(key, val string) RestOption {
	return func(rc *RestClient) {
		rc.headers.Add(key, val)
	}
}
