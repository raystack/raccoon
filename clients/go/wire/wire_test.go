package wire

import (
	"testing"

	pb "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/raccoon/v1beta1"

	"github.com/stretchr/testify/assert"
)

func TestJsonWireMarshal(t *testing.T) {
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

	wjson := &JsonWire{}
	buf, err := wjson.Marshal(msg)
	if err != nil {
		t.Errorf("wireJson.Marshal (%v) failed with %v; want success", msg, err)
	}

	got := &pb.SendEventRequest{}
	if err := wjson.Unmarshal(buf, got); err != nil {
		t.Errorf("wireJson.Unmarshal (%v) failed with %v; want success", buf, err)
	}

	assert.Equal("application/json", wjson.ContentType())
	assert.Equal(msg.ReqGuid, got.ReqGuid)
	assert.Equal(len(msg.Events), len(got.Events))
	assert.Equal(msg.Events[0].Type, got.Events[0].Type)
	assert.Equal(msg.Events[0].EventBytes, got.Events[0].EventBytes)
}

func TestProtoWireMarshal(t *testing.T) {
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

	wproto := &ProtoWire{}
	buf, err := wproto.Marshal(msg)
	if err != nil {
		t.Errorf("wireProto.Marshal (%v) failed with %v; want success", msg, err)
	}

	got := &pb.SendEventRequest{}
	if err := wproto.Unmarshal(buf, got); err != nil {
		t.Errorf("wireProto.Unmarshal (%v) failed with %v; want success", buf, err)
	}

	assert.Equal("application/proto", wproto.ContentType())
	assert.Equal(msg.ReqGuid, got.ReqGuid)
	assert.Equal(len(msg.Events), len(got.Events))
	assert.Equal(msg.Events[0].Type, got.Events[0].Type)
	assert.Equal(msg.Events[0].EventBytes, got.Events[0].EventBytes)
}
