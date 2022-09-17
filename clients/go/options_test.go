package raccoon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestOptionsSet(t *testing.T) {
	assert := assert.New(t)

	url := "http://localhost:8080"
	key := "authorization"
	val := "123"

	rc := NewRest(
		WithUrl(url),
		WithHeader(key, val),
		WithMarshaler(JSON))

	assert.NotNil(rc.Marshal)
	assert.Equal(url, rc.Url)
	assert.Equal(val, rc.headers.Get(key))
}
