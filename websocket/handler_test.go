package websocket

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"net/http"
	"net/http/httptest"
	"source.golabs.io/mobile/clickstream-go-proto/gojek/clickstream/de"
	"strings"
	"testing"
	"time"

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

func TestHandlerWSEvents(t *testing.T) {
	// ---- Setup ----
	hlr := &Handler{

		websocketUpgrader: websocket.Upgrader{
			ReadBufferSize:  10240,
			WriteBufferSize: 10240,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		bufferChannel: make(chan []*de.CSEventMessage, 10),
	}
	ts := httptest.NewServer(Router(hlr))
	defer ts.Close()

	url := "ws" + strings.TrimPrefix(ts.URL+"/api/v1/events", "http")
	header := http.Header{
		"User-ID": []string{"test1-user1"},
	}

	t.Run("Should return success response after successfully push to channel", func(t *testing.T) {
		wss, _, err := websocket.DefaultDialer.Dial(url, header)
		require.NoError(t, err)

		request := &de.EventRequest{
			ReqGuid:  "1234",
			SentTime: time.Now().Unix(),
			Data:     nil,
		}
		serializedRequest, _ := proto.Marshal(request)

		err = wss.WriteMessage(websocket.BinaryMessage, serializedRequest)
		require.NoError(t, err)

		responseMsgType, response, err := wss.ReadMessage()
		require.NoError(t, err)

		resp := &de.EventResponse{}
		proto.Unmarshal(response, resp)
		assert.Equal(t, responseMsgType, websocket.BinaryMessage)
		assert.Equal(t, request.ReqGuid, resp.GetData()["req_guid"])
		assert.Equal(t, de.Status_SUCCESS, resp.GetStatus())
		assert.Equal(t, de.Code_OK, resp.GetCode())
		assert.Equal(t, "", resp.GetReason())
	})

	t.Run("Should return unknown request when request fail to deserialize", func(t *testing.T) {
		wss, _, err := websocket.DefaultDialer.Dial(url, header)
		require.NoError(t, err)

		err = wss.WriteMessage(websocket.BinaryMessage, []byte{1, 2, 3, 4, 1, 2})
		require.NoError(t, err)

		responseMsgType, response, err := wss.ReadMessage()
		require.NoError(t, err)

		resp := &de.EventResponse{}
		proto.Unmarshal(response, resp)
		assert.Equal(t, responseMsgType, websocket.BinaryMessage)
		assert.Equal(t, de.Status_ERROR, resp.GetStatus())
		assert.Equal(t, de.Code_BAD_REQUEST, resp.GetCode())
		assert.Empty(t, resp.GetData())
	})
}
