package rest

import (
	"testing"

	"github.com/raystack/raccoon/clients/go/serializer"
	"github.com/stretchr/testify/assert"
)

func TestRestOptionsSet(t *testing.T) {
	assert := assert.New(t)

	url := "http://localhost:8080"
	key := "authorization"
	val := "123"

	rc, err := New(
		WithUrl(url),
		WithHeader(key, val),
		WithSerializer(serializer.JSON))

	assert.NoError(err)
	assert.NotNil(rc.serialize)
	assert.Equal(url, rc.url)
	assert.Equal(val, rc.headers.Get(key))
}
