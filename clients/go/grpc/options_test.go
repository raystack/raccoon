package grpc

import (
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/stretchr/testify/assert"
)

func TestGrpcOptionsSet(t *testing.T) {
	assert := assert.New(t)

	addr := "localhost:8080"
	key := "x-UniqueId"
	val := "123"
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	gc, err := NewGrpc(
		WithAddr(addr),
		WithHeader(key, val),
		WithDialOptions(opts...))

	assert.NoError(err)
	assert.Equal(addr, gc.addr)
	assert.Equal(1, len(gc.dialOptions))
	assert.Equal("123", gc.headers[key])
}

func TestGrpcDialOptionsSet(t *testing.T) {
	assert := assert.New(t)

	addr := "localhost:8080"
	key := "authorization"
	val := "123"
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	gc, err := NewGrpc(
		WithAddr(addr),
		WithHeader(key, val),
		WithDialOptions(opts...),
	)

	assert.NoError(err)
	assert.Equal(addr, gc.addr)
	assert.Equal(len(opts), len(gc.dialOptions))
}
