package serializer

import (
	"testing"

	"encoding/json"

	pb "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/raccoon/v1beta1"
	"google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"
)

func TestJsonSerializer(t *testing.T) {
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
		t.Errorf("json.Serializer (%v) failed with %v; want success", msg, err)
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

func TestProtoSerializer(t *testing.T) {
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
		t.Errorf("proto.Serializer (%v) failed with %v; want success", msg, err)
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
