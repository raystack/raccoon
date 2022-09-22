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
	Url        string
	Serialize  serializer.SerializerFunc
	Wire       wire.WireMarshaler
	httpclient heimdall.Client
	headers    http.Header
}

// RestOption represents the client options.
type RestOption func(*RestClient)

// WithUrl sets the service address
func WithUrl(url string) RestOption {
	return func(rc *RestClient) {
		rc.Url = url
	}
}

// WithSerializer sets the serializer for the raccoon message.
func WithSerializer(s serializer.SerializerFunc) RestOption {
	return func(rc *RestClient) {
		rc.Serialize = s
	}
}

// WithHeader sets the http header for the request.
func WithHeader(key, val string) RestOption {
	return func(rc *RestClient) {
		rc.headers.Add(key, val)
	}
}
