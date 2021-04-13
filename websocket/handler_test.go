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
	hlr := &Handler{
		websocketUpgrader: websocket.Upgrader{
			ReadBufferSize:  10240,
			WriteBufferSize: 10240,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		user:              NewUserStore(2),
		bufferChannel:     make(chan EventsBatch, 10),
		PongWaitInterval:  time.Duration(60 * time.Second),
		WriteWaitInterval: time.Duration(5 * time.Second),
		PingChannel:       make(chan connection, 100),
		UniqConnIDHeader:      "x-user-id",
	}
	ts := httptest.NewServer(Router(hlr))
	defer ts.Close()

	url := "ws" + strings.TrimPrefix(ts.URL+"/api/v1/events", "http")
	header := http.Header{
		"x-user-id": []string{"test1-user1"},
	}

	t.Run("Should return success response after successfully push to channel", func(t *testing.T) {
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
		wss, _, err := websocket.DefaultDialer.Dial(url, http.Header{
			"x-user-id": []string{"test2-user2"},
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

	t.Run("Should close subsequence connection of the same user", func(t *testing.T) {
		ts := httptest.NewServer(Router(hlr))
		defer ts.Close()
		url := "ws" + strings.TrimPrefix(ts.URL+"/api/v1/events", "http")
		header := http.Header{
			"x-user-id": []string{"test1-user1"},
		}
		firstWss, _, err := websocket.DefaultDialer.Dial(url, header)
		require.NoError(t, err)

		secondWss, _, err := websocket.DefaultDialer.Dial(url, header)
		require.NoError(t, err)
		_, message, err := secondWss.ReadMessage()
		p := &pb.EventResponse{}
		proto.Unmarshal(message, p)
		assert.Equal(t, p.Code, pb.Code_MAX_USER_LIMIT_REACHED)
		assert.Equal(t, p.Status, pb.Status_ERROR)
		_, _, err = secondWss.ReadMessage()
		assert.True(t, websocket.IsCloseError(err, websocket.ClosePolicyViolation))
		assert.Equal(t, "Duplicate connection", err.(*websocket.CloseError).Text)
		firstWss.Close()
	})

	t.Run("Should close new connection when reach max connection", func(t *testing.T) {
		ts := httptest.NewServer(Router(hlr))
		defer ts.Close()
		url := "ws" + strings.TrimPrefix(ts.URL+"/api/v1/events", "http")
		header := http.Header{
			"x-user-id": []string{"test1-user1"},
		}
		websocket.DefaultDialer.Dial(url, http.Header{"x-user-id": []string{"test1-user2"}})
		websocket.DefaultDialer.Dial(url, http.Header{"x-user-id": []string{"test1-user3"}})

		ws, _, err := websocket.DefaultDialer.Dial(url, header)
		require.NoError(t, err)
		_, message, err := ws.ReadMessage()
		p := &pb.EventResponse{}
		proto.Unmarshal(message, p)
		assert.Equal(t, p.Code, pb.Code_MAX_CONNECTION_LIMIT_REACHED)
		assert.Equal(t, p.Status, pb.Status_ERROR)
		_, _, err = ws.ReadMessage()
		assert.True(t, websocket.IsCloseError(err, websocket.ClosePolicyViolation))
		assert.Equal(t, "Max connection reached", err.(*websocket.CloseError).Text)
	})

	t.Run("Should decrement total connection when client close the conn", func(t *testing.T) {
		ts := httptest.NewServer(Router(hlr))
		defer ts.Close()
		url := "ws" + strings.TrimPrefix(ts.URL+"/api/v1/events", "http")
		header := http.Header{
			"x-user-id": []string{"test1-user1"},
		}
		firstWs, _, _ := websocket.DefaultDialer.Dial(url, http.Header{"x-user-id": []string{"test1-user2"}})
		firstWs.Close()
		websocket.DefaultDialer.Dial(url, http.Header{"x-user-id": []string{"test1-user3"}})

		_, _, err := websocket.DefaultDialer.Dial(url, header)
		assert.Equal(t, 2, hlr.user.TotalUsers())
		assert.Empty(t, err)
	})
}
