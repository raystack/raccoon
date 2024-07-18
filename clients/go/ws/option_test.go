package ws

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func noopHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	c.Close()
}

func TestWebSocketOptionsSet(t *testing.T) {
	assert := assert.New(t)

	srv := httptest.NewServer(http.HandlerFunc(noopHandler))
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	assert.NoError(err)

	url := fmt.Sprintf("ws://%s/api/v1/events", u.Host)
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
