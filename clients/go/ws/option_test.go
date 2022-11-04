package ws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebSocketOptionsSet(t *testing.T) {
	assert := assert.New(t)

	url := "ws://localhost:8080/api/v1/events"
	key := "authorization"
	val := "123"

	ws, err := New(
		WithUrl(url),
		WithHeader(key, val))

	assert.NoError(err)
	assert.NotNil(ws.serialize)
	assert.Equal(url, ws.url)
	assert.Equal(val, ws.headers.Get(key))
}
