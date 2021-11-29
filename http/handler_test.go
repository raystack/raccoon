package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"raccoon/collection"
	"raccoon/http/grpc"
	"raccoon/http/rest"
	"raccoon/http/websocket"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPingHandler(t *testing.T) {
	hlr := &Handler{
		wh: &websocket.Handler{},
		rh: rest.NewHandler(),
		gh: &grpc.Handler{},
	}
	collector := new(collection.MockCollector)
	collector.On("Collect", mock.Anything, mock.Anything).Return(nil)
	ts := httptest.NewServer(Router(hlr, collector))
	defer ts.Close()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/ping", ts.URL), nil)

	httpClient := http.Client{}
	res, _ := httpClient.Do(req)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}
