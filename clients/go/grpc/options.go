package grpc

import (
	"google.golang.org/grpc"

	raccoon "github.com/odpf/raccoon/clients/go"
	"github.com/odpf/raccoon/clients/go/serializer"
)

// GrpcClient is the http implementation
type GrpcClient struct {
	raccoon.Client
	addr        string
	client      *grpc.ClientConn
	serialize   serializer.SerializerFunc
	headers     map[string]string
	dialOptions []grpc.DialOption
}

// GrpcOption represents the client options.
type GrpcOption func(*GrpcClient)

// WithAddr sets the service address
func WithAddr(addr string) GrpcOption {
	return func(gc *GrpcClient) {
		gc.addr = addr
	}
}

// WithSerializer sets the serializer for the raccoon message.
func WithSerializer(s serializer.SerializerFunc) GrpcOption {
	return func(gc *GrpcClient) {
		gc.serialize = s
	}
}

// WithHeader sets the grpc metadata for the request.
func WithHeader(key, val string) GrpcOption {
	return func(gc *GrpcClient) {
		gc.headers[key] = val
	}
}

// WithDialOptions sets the grpc dial options.
func WithDialOptions(opts ...grpc.DialOption) GrpcOption {
	return func(gc *GrpcClient) {
		gc.dialOptions = opts
	}
}
