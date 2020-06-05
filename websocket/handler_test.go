package websocket

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPingHandler(t *testing.T) {
	ts := httptest.NewServer(Router(nil))
	defer ts.Close()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/ping", ts.URL), nil)

	httpClient := http.Client{}
	res, _ := httpClient.Do(req)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestWSHandlerSendsAcknowledgement(t *testing.T) {

	hlr := &Handler{
		websocketUpgrader: websocket.Upgrader{
			ReadBufferSize:  10240,
			WriteBufferSize: 10240,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	ts := httptest.NewServer(Router(hlr))
	defer ts.Close()
	
	url := "ws" + strings.TrimPrefix(ts.URL + "/api/v1/events", "http")
	header := http.Header{
		"User-ID":	[]string{"test1-user1"},
	}
	log.Println(fmt.Sprintf("%s", ts.URL))
	wss, _, err := websocket.DefaultDialer.Dial(url, header)
	require.NoError(t, err)

	err = wss.WriteMessage(websocket.BinaryMessage, []byte("TestWsServerEndPoint"))
	require.NoError(t, err)

	responseMsgType, response, err := wss.ReadMessage()
	require.NoError(t, err)

	assert.Equal(t, responseMsgType, websocket.TextMessage)
	assert.Equal(t, "batch-id: test1-user1", string(response))
}
