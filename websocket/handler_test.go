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
		ConnIDHeader:      "x-user-id",
		ConnTypeHeader:    "",
	}
	ts := httptest.NewServer(Router(hlr))
	defer ts.Close()

	url := "ws" + strings.TrimPrefix(ts.URL+"/api/v1/events", "http")
	header := http.Header{
		"x-user-id": []string{"test1-user1"},
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
		w1, _, err := websocket.DefaultDialer.Dial(url, header)
		defer w1.Close()
		require.NoError(t, err)

		w2, _, err := websocket.DefaultDialer.Dial(url, header)
		defer w2.Close()
		require.NoError(t, err)
		_, message, err := w2.ReadMessage()
		p := &pb.EventResponse{}
		proto.Unmarshal(message, p)
		assert.Equal(t, p.Code, pb.Code_MAX_USER_LIMIT_REACHED)
		assert.Equal(t, p.Status, pb.Status_ERROR)
		_, _, err = w2.ReadMessage()
		assert.True(t, websocket.IsCloseError(err, websocket.ClosePolicyViolation))
		assert.Equal(t, "Duplicate connection", err.(*websocket.CloseError).Text)
	})

	t.Run("Should accept connection with same id and different type", func(t *testing.T) {
		hlr.ConnTypeHeader = "test-type"
		ts := httptest.NewServer(Router(hlr))
		defer func() { hlr.ConnTypeHeader = "" }()
		defer ts.Close()

		url := "ws" + strings.TrimPrefix(ts.URL+"/api/v1/events", "http")
		header := http.Header{
			"x-user-id": []string{"test1-user1"},
		}
		header["test-type"] = []string{"type-1"}
		w1, _, err := websocket.DefaultDialer.Dial(url, header)
		defer w1.Close()
		require.NoError(t, err)

		header["test-type"] = []string{"type-2"}
		w2, _, err := websocket.DefaultDialer.Dial(url, header)
		defer w2.Close()

		request := &pb.EventRequest{
			ReqGuid:  "1234",
			SentTime: ptypes.TimestampNow(),
			Events:   nil,
		}
		serializedRequest, _ := proto.Marshal(request)

		err = w2.WriteMessage(websocket.BinaryMessage, serializedRequest)

		require.NoError(t, err)
		_, message, err := w2.ReadMessage()
		p := &pb.EventResponse{}
		proto.Unmarshal(message, p)
		assert.Equal(t, pb.Code_OK, p.Code)
		assert.Equal(t, pb.Status_SUCCESS, p.Status)
	})

	t.Run("Should close new connection when reach max connection", func(t *testing.T) {
		ts := httptest.NewServer(Router(hlr))
		defer ts.Close()

		url := "ws" + strings.TrimPrefix(ts.URL+"/api/v1/events", "http")
		header := http.Header{
			"x-user-id": []string{"test1-user1"},
		}
		w1, _, _ := websocket.DefaultDialer.Dial(url, http.Header{"x-user-id": []string{"test1-user2"}})
		defer w1.Close()
		w2, _, _ := websocket.DefaultDialer.Dial(url, http.Header{"x-user-id": []string{"test1-user3"}})
		defer w2.Close()

		w3, _, err := websocket.DefaultDialer.Dial(url, header)
		defer w3.Close()
		require.NoError(t, err)
		_, message, err := w3.ReadMessage()
		p := &pb.EventResponse{}
		proto.Unmarshal(message, p)
		assert.Equal(t, p.Code, pb.Code_MAX_CONNECTION_LIMIT_REACHED)
		assert.Equal(t, p.Status, pb.Status_ERROR)
		_, _, err = w3.ReadMessage()
		assert.True(t, websocket.IsCloseError(err, websocket.ClosePolicyViolation))
		assert.Equal(t, "Max connection reached", err.(*websocket.CloseError).Text)
	})

	t.Run("Should decrement total connection when client close the conn", func(t *testing.T) {
		ts := httptest.NewServer(Router(hlr))
		defer ts.Close()

		url := "ws" + strings.TrimPrefix(ts.URL+"/api/v1/events", "http")
		w1, _, _ := websocket.DefaultDialer.Dial(url, http.Header{"x-user-id": []string{"test1-user2"}})
		defer w1.Close()
		w2, _, _ := websocket.DefaultDialer.Dial(url, http.Header{"x-user-id": []string{"test1-user3"}})
		defer w2.Close()
		w3, _, err := websocket.DefaultDialer.Dial(url, http.Header{"x-user-id": []string{"test1-user1"}})
		defer w3.Close()

		assert.Equal(t, 2, hlr.user.TotalUsers())
		assert.Empty(t, err)
	})
}
