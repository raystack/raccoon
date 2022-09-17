package raccoon

import (
	"net/http"

	"github.com/gojek/heimdall/v7"
)

// RestClient is the http implementation
type RestClient struct {
	Url        string
	Marshal    MarshalFunc
	Wire       WireMarshaler
	httpclient heimdall.Client
	headers    http.Header
}

// RestOption represents the client options.
type RestOption func(*RestClient)

// WithUrl sets the service address
func WithUrl(url string) RestOption {
	return func(o *RestClient) {
		o.Url = url
	}
}

// WithMarshaler sets the marshaler for the event message.
func WithMarshaler(m MarshalFunc) RestOption {
	return func(o *RestClient) {
		o.Marshal = m
	}
}

// WithHeader sets the http header for the request.
func WithHeader(key, val string) RestOption {
	return func(o *RestClient) {
		o.headers.Add(key, val)
	}
}
