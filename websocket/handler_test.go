package websocket

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"raccoon/logger"
	"raccoon/metrics"
	"strings"
	"testing"
	"time"

	"raccoon/websocket/connection"
	pb "raccoon/websocket/proto"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type void struct{}

func (v void) Write(_ []byte) (int, error) {
	return 0, nil
}
func TestMain(t *testing.M) {
	logger.SetOutput(void{})
	metrics.SetVoid()
	os.Exit(t.Run())
}

func TestPingHandler(t *testing.T) {
	ts := httptest.NewServer(Router(nil))
	defer ts.Close()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/ping", ts.URL), nil)

	httpClient := http.Client{}
	res, _ := httpClient.Do(req)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestHandler_HandlerWSEvents(t *testing.T) {
	// ---- Setup ----
	upgrader := connection.NewUpgrader(connection.UpgraderConfig{
		ReadBufferSize:    10240,
		WriteBufferSize:   10240,
		CheckOrigin:       false,
		MaxUser:           2,
		PongWaitInterval:  time.Duration(60 * time.Second),
		WriteWaitInterval: time.Duration(5 * time.Second),
		ConnIDHeader:      "X-User-ID",
		ConnGroupHeader:   "string",
	})
	hlr := &Handler{
		upgrader:      upgrader,
		bufferChannel: make(chan EventsBatch, 10),
		PingChannel:   make(chan connection.Conn, 100),
	}
	ts := httptest.NewServer(Router(hlr))
	defer ts.Close()

	url := "ws" + strings.TrimPrefix(ts.URL+"/api/v1/events", "http")
	header := http.Header{
		"X-User-ID": []string{"test1-user1"},
	}

	t.Run("Should return success response after successfully push to channel", func(t *testing.T) {
		ts = httptest.NewServer(Router(hlr))
		defer ts.Close()

		wss, _, err := websocket.DefaultDialer.Dial(url, header)
		require.NoError(t, err)

		request := &pb.EventRequest{
			ReqGuid:  "1234",
			SentTime: ptypes.TimestampNow(),
			Events:   nil,
		}
		serializedRequest, _ := proto.Marshal(request)

		err = wss.WriteMessage(websocket.BinaryMessage, serializedRequest)
		defer wss.Close()
		require.NoError(t, err)

		responseMsgType, response, err := wss.ReadMessage()
		require.NoError(t, err)

		resp := &pb.EventResponse{}
		proto.Unmarshal(response, resp)
		assert.Equal(t, responseMsgType, websocket.BinaryMessage)
		assert.Equal(t, request.ReqGuid, resp.GetData()["req_guid"])
		assert.Equal(t, pb.Status_SUCCESS, resp.GetStatus())
		assert.Equal(t, pb.Code_OK, resp.GetCode())
		assert.Equal(t, "", resp.GetReason())
	})

	t.Run("Should return unknown request when request fail to deserialize", func(t *testing.T) {
		ts = httptest.NewServer(Router(hlr))
		defer ts.Close()

		wss, _, err := websocket.DefaultDialer.Dial(url, http.Header{
			"X-User-ID": []string{"test2-user2"},
		})
		defer wss.Close()
		require.NoError(t, err)

		err = wss.WriteMessage(websocket.BinaryMessage, []byte{1, 2, 3, 4, 1, 2})
		require.NoError(t, err)

		responseMsgType, response, err := wss.ReadMessage()
		require.NoError(t, err)

		resp := &pb.EventResponse{}
		proto.Unmarshal(response, resp)
		assert.Equal(t, responseMsgType, websocket.BinaryMessage)
		assert.Equal(t, pb.Status_ERROR, resp.GetStatus())
		assert.Equal(t, pb.Code_BAD_REQUEST, resp.GetCode())
		assert.Empty(t, resp.GetData())
	})
}
