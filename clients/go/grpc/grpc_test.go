package grpc

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"
	"time"

	pb "go.buf.build/odpf/gw/odpf/proton/odpf/raccoon/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	raccoon "github.com/odpf/raccoon/clients/go"
	"github.com/stretchr/testify/assert"
)

const connId string = "X-UniqueId"

type mockEventServiceServer struct {
	pb.UnimplementedEventServiceServer
}

func (*mockEventServiceServer) SendEvent(ctx context.Context, req *pb.SendEventRequest) (*pb.SendEventResponse, error) {
	metadata, _ := metadata.FromIncomingContext(ctx)
	id := metadata.Get(connId)
	if len(id) == 0 {
		return nil, errors.New("conn id header is missing")
	}

	pe := &PageEvent{}
	if err := proto.Unmarshal(req.Events[0].EventBytes, pe); err != nil {
		return nil, err
	}

	return &pb.SendEventResponse{
		Status:   pb.Status_STATUS_SUCCESS,
		Code:     pb.Code_CODE_OK,
		SentTime: time.Now().Unix(),
		Data: map[string]string{
			"req_guid": req.ReqGuid,
		},
	}, nil
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	pb.RegisterEventServiceServer(server, &mockEventServiceServer{})
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()
	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestGrpcClientSend(t *testing.T) {
	assert := assert.New(t)

	gc, err := NewGrpc(
		WithAddr(""),
		WithHeader(connId, "123"),
		WithDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(dialer()),
			grpc.WithBlock()),
	)
	assert.NoError(err)
	defer gc.Close()

	reqGuid, resp, err := gc.Send([]*raccoon.Event{
		{
			Type: "page",
			Data: &PageEvent{
				EventGuid: "guid-123",
				EventName: "page",
				SentTime:  timestamppb.Now(),
			},
		},
	})

	assert.NotEmpty(reqGuid)
	assert.Nil(err)
	assert.Equal(int32(1), resp.Status)
	assert.NotNil(resp.SentTime)
}
