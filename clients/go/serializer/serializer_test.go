package serializer

import (
	"testing"

	"encoding/json"

	pb "go.buf.build/odpf/gw/odpf/proton/odpf/raccoon/v1beta1"
	"google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"
)

func TestJsonMarshal(t *testing.T) {
	assert := assert.New(t)

	msg := &pb.SendEventRequest{
		ReqGuid: "reqId",
		Events: []*pb.Event{
			{
				Type:       "page",
				EventBytes: []byte(`{"id": "foo"}`),
			},
		},
	}

	buf, err := JSON(msg)
	if err != nil {
		t.Errorf("json.Marshal (%v) failed with %v; want success", msg, err)
	}

	got := &pb.SendEventRequest{}
	if err := json.Unmarshal(buf, got); err != nil {
		t.Errorf("json.Unmarshal (%v) failed with %v; want success", buf, err)
	}

	assert.Equal(msg.ReqGuid, got.ReqGuid)
	assert.Equal(len(msg.Events), len(got.Events))
	assert.Equal(msg.Events[0].Type, got.Events[0].Type)
	assert.Equal(msg.Events[0].EventBytes, got.Events[0].EventBytes)
}

func TestProtoMarshal(t *testing.T) {
	assert := assert.New(t)

	msg := &pb.SendEventRequest{
		ReqGuid: "reqId",
		Events: []*pb.Event{
			{
				Type:       "page",
				EventBytes: []byte(`{"id": "foo"}`),
			},
		},
	}

	buf, err := PROTO(msg)
	if err != nil {
		t.Errorf("proto.Marshal (%v) failed with %v; want success", msg, err)
	}

	got := &pb.SendEventRequest{}
	if err := proto.Unmarshal(buf, got); err != nil {
		t.Errorf("proto.Unmarshal (%v) failed with %v; want success", buf, err)
	}

	assert.Equal(msg.ReqGuid, got.ReqGuid)
	assert.Equal(len(msg.Events), len(got.Events))
	assert.Equal(msg.Events[0].Type, got.Events[0].Type)
	assert.Equal(msg.Events[0].EventBytes, got.Events[0].EventBytes)
}
