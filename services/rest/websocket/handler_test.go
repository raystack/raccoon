package websocket

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	pb "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/raccoon/v1beta1"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/goto/raccoon/collection"
	"github.com/goto/raccoon/logger"
	"github.com/goto/raccoon/metrics"
	"github.com/goto/raccoon/services/rest/websocket/connection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type void struct{}

func (v void) Write(_ []byte) (int, error) {
	return 0, nil
}
func TestMain(t *testing.M) {
	logger.SetOutput(void{})
	metrics.SetVoid()
}

func TestNewHandler(t *testing.T) {
	type args struct {
		pingC chan connection.Conn
	}

	ugConfig := connection.UpgraderConfig{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		CheckOrigin:       false,
		MaxUser:           100,
		PongWaitInterval:  60,
		WriteWaitInterval: 60,
		ConnIDHeader:      "x-conn-id",
		ConnGroupHeader:   "x-group",
	}
	pingC := make(chan connection.Conn)
	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{
			name: "creating a new handler",
			args: args{
				pingC: pingC,
			},
			want: &Handler{
				upgrader:    connection.NewUpgrader(ugConfig),
				PingChannel: pingC,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandler(tt.args.pingC, &collection.MockCollector{}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_Table(t *testing.T) {
	table := &connection.Table{}
	type fields struct {
		upgrader    *connection.Upgrader
		PingChannel chan connection.Conn
	}
	tests := []struct {
		name   string
		fields fields
		want   *connection.Table
	}{
		{
			name: "return table",
			fields: fields{
				upgrader: &connection.Upgrader{
					Table: table,
				},
			},
			want: table,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				upgrader:    tt.fields.upgrader,
				PingChannel: tt.fields.PingChannel,
			}
			if got := h.Table(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.Table() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_GETHandlerWSEvents(t *testing.T) {
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
		upgrader: upgrader,

		PingChannel: make(chan connection.Conn, 100),
	}
	ts := httptest.NewServer(getRouter(hlr))
	defer ts.Close()

	url := "ws" + strings.TrimPrefix(ts.URL+"/api/v1/events", "http")
	header := http.Header{
		"X-User-ID": []string{"test1-user1"},
	}

	t.Run("Should return success response after successfully push to channel", func(t *testing.T) {
		ts = httptest.NewServer(getRouter(hlr))
		defer ts.Close()

		wss, _, err := websocket.DefaultDialer.Dial(url, header)
		require.NoError(t, err)

		request := &pb.SendEventRequest{
			ReqGuid:  "1234",
			SentTime: timestamppb.Now(),
			Events:   nil,
		}
		serializedRequest, _ := proto.Marshal(request)

		err = wss.WriteMessage(websocket.BinaryMessage, serializedRequest)
		defer wss.Close()
		require.NoError(t, err)

		responseMsgType, response, err := wss.ReadMessage()
		require.NoError(t, err)

		resp := &pb.SendEventResponse{}
		proto.Unmarshal(response, resp)
		assert.Equal(t, responseMsgType, websocket.BinaryMessage)
		assert.Equal(t, request.ReqGuid, resp.GetData()["req_guid"])
		assert.Equal(t, pb.Status_STATUS_SUCCESS, resp.GetStatus())
		assert.Equal(t, pb.Code_CODE_OK, resp.GetCode())
		assert.Equal(t, "", resp.GetReason())
	})

	t.Run("Should return unknown request when request fail to deserialize", func(t *testing.T) {
		ts = httptest.NewServer(getRouter(hlr))
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

		resp := &pb.SendEventResponse{}
		proto.Unmarshal(response, resp)
		assert.Equal(t, responseMsgType, websocket.BinaryMessage)
		assert.Equal(t, pb.Status_STATUS_ERROR, resp.GetStatus())
		assert.Equal(t, pb.Code_CODE_BAD_REQUEST, resp.GetCode())
		assert.Empty(t, resp.GetData())
	})
}

func getRouter(hlr *Handler) http.Handler {
	collector := new(collection.MockCollector)
	collector.On("Collect", mock.Anything, mock.Anything).Return(nil)
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()
	subRouter.HandleFunc("/events", hlr.HandlerWSEvents).Methods(http.MethodGet).Name("events")
	return router
}
