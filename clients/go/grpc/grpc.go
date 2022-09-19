package grpc

import (
	"context"

	pb "go.buf.build/odpf/gw/odpf/proton/odpf/raccoon/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	raccoon "github.com/odpf/raccoon/clients/go"
	"github.com/odpf/raccoon/clients/go/serializer"
)

// NewGrpc creates the new grpc client with provided options.
func NewGrpc(options ...GrpcOption) (*GrpcClient, error) {
	gc := &GrpcClient{
		Serialize:   serializer.PROTO,
		headers:     make(map[string]string),
		dialOptions: []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock()},
	}

	for _, opt := range options {
		opt(gc)
	}

	client, err := grpc.Dial(gc.Addr, gc.dialOptions...)
	if err != nil {
		return nil, err
	}
	gc.client = client
	return gc, err
}

// Send sends the events to the raccoon service
func (g *GrpcClient) Send(events []*raccoon.Event) (string, *raccoon.Response, error) {
	reqId := uuid.NewString()

	e := []*pb.Event{}
	for _, ev := range events {
		b, err := g.Serialize(ev.Data)
		if err != nil {
			return reqId, nil, err
		}
		e = append(e, &pb.Event{
			EventBytes: b,
			Type:       ev.Type,
		})
	}

	svc := pb.NewEventServiceClient(g.client)
	meta := metadata.New(g.headers)
	resp, err := svc.SendEvent(metadata.NewOutgoingContext(context.Background(), meta), &pb.SendEventRequest{
		ReqGuid:  reqId,
		Events:   e,
		SentTime: timestamppb.Now(),
	})
	if err != nil {
		return reqId, nil, err
	}

	return reqId, &raccoon.Response{
		Status:   int32(resp.Status),
		Code:     int32(resp.Code),
		SentTime: resp.SentTime,
		Data:     resp.Data,
	}, nil
}

// Close closes the grpc connection.
func (g *GrpcClient) Close() {
	g.client.Close()
}
