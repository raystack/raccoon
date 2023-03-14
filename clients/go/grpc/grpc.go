package grpc

import (
	"context"
	"fmt"

	"github.com/goto/raccoon/clients/go/log"

	pbgrpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/raccoon/v1beta1/raccoonv1beta1grpc"
	pb "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/raccoon/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	raccoon "github.com/goto/raccoon/clients/go"
	"github.com/goto/raccoon/clients/go/retry"
	"github.com/goto/raccoon/clients/go/serializer"
)

// New creates the new grpc client with provided options.
func New(options ...Option) (*Grpc, error) {
	gc := &Grpc{
		serialize:   serializer.PROTO,
		headers:     make(map[string]string),
		dialOptions: []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock()},
		retryMax:    retry.DefaultRetryMax,
		retryWait:   retry.DefaultRetryWait,
		logger:      log.Default(),
	}

	for _, opt := range options {
		opt(gc)
	}

	client, err := grpc.Dial(gc.addr, gc.dialOptions...)
	if err != nil {
		return nil, err
	}
	gc.client = client
	return gc, err
}

// Send sends the events to the raccoon service
func (gc *Grpc) Send(events []*raccoon.Event) (string, *raccoon.Response, error) {
	reqId := uuid.NewString()
	gc.logger.Infof("started request, addr: %s, req-id: %s", gc.addr, reqId)
	defer gc.logger.Infof("ended request, addr: %s, req-id: %s", gc.addr, reqId)

	e := []*pb.Event{}
	for _, ev := range events {
		b, err := gc.serialize(ev.Data)
		if err != nil {
			gc.logger.Errorf("serialize, addr: %s, req-id: %s, %+v", gc.addr, reqId, err)
			return reqId, nil, err
		}
		e = append(e, &pb.Event{
			EventBytes: b,
			Type:       ev.Type,
		})
	}

	svc := pbgrpc.NewEventServiceClient(gc.client)
	meta := metadata.New(gc.headers)
	racReq := &pb.SendEventRequest{
		ReqGuid:  reqId,
		Events:   e,
		SentTime: timestamppb.Now(),
	}

	var resp *pb.SendEventResponse
	err := retry.Do(gc.retryWait, gc.retryMax, func() error {
		res, err := svc.SendEvent(metadata.NewOutgoingContext(context.Background(), meta), racReq)
		if err != nil {
			return err
		}

		if res.Status != pb.Status_STATUS_SUCCESS {
			return fmt.Errorf("error from raccoon addr: %s, req-id: %s, status: %d, code: %d, data: %+v", gc.addr, reqId, res.Status, res.Code, res.Data)
		}

		resp = res
		return nil
	})
	if err != nil {
		gc.logger.Errorf("send, addr: %s, req-id: %s, %+v", gc.addr, reqId, err)
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
func (gc *Grpc) Close() {
	gc.client.Close()
}
