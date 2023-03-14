package ws

import (
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	pb "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/raccoon/v1beta1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	raccoon "github.com/goto/raccoon/clients/go"
	"github.com/goto/raccoon/clients/go/serializer"
	"github.com/goto/raccoon/clients/go/testdata"
	"github.com/stretchr/testify/assert"
)

var upgrader = websocket.Upgrader{}

func Handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		_, rMsg, err := c.ReadMessage()
		if err != nil {
			break
		}

		req := &pb.SendEventRequest{}
		err = proto.Unmarshal(rMsg, req)
		if err != nil {
			break
		}

		msg, err := serializer.PROTO(&pb.SendEventResponse{
			Status:   1,
			Code:     1,
			Reason:   "",
			SentTime: timestamppb.Now().AsTime().Unix(),
			Data: map[string]string{
				"req_guid": req.ReqGuid,
			},
		})
		if err != nil {
			break
		}
		err = c.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			break
		}
	}
}

func TestWebSocket(t *testing.T) {
	assert := assert.New(t)
	svr := httptest.NewServer(http.HandlerFunc(Handler))
	defer svr.Close()

	url := "ws" + strings.TrimPrefix(svr.URL, "http")

	client, err := New(
		WithUrl(url),
		WithHeader("x-user-id", "123"),
		WithHeader("x-user-type", "gojek"))
	assert.NoError(err)

	defer client.Close()

	reqGuid, err := client.Send([]*raccoon.Event{
		{
			Type: "page",
			Data: &testdata.PageEvent{
				EventGuid: uuid.NewString(),
				EventName: "clicked",
				SentTime:  timestamppb.Now(),
			},
		},
	})
	assert.NotEmpty(reqGuid)
	assert.NoError(err)

	resp := <-client.EventAcks()
	assert.NotNil(resp)
	assert.Equal(int32(1), resp.Status)
	assert.Equal(int32(1), resp.Code)
	assert.Equal(reqGuid, resp.Data["req_guid"])
}
