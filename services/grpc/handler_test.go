package grpc

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/raystack/raccoon/collector"
	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/logger"
	"github.com/raystack/raccoon/metrics"
	pb "github.com/raystack/raccoon/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestHandler_SendEvent(t *testing.T) {
	// todo(turtledev): refactor this test
	// it has the following issues:
	//  1. collector is shared across all tests cases
	//  2. test case specific parameters are kept in shared scope

	type fields struct {
		C   collector.Collector
		ack config.AckType
	}
	type args struct {
		ctx context.Context
		req *pb.SendEventRequest
	}

	logger.SetOutput(io.Discard)
	metrics.SetVoid()
	ctx := context.Background()
	meta := metadata.MD{}
	meta.Set(config.Server.Websocket.Conn.GroupHeader, "group")
	meta.Set(config.Server.Websocket.Conn.IDHeader, "1235")
	sentTime := timestamppb.Now()
	req := &pb.SendEventRequest{
		ReqGuid:  "abcd",
		SentTime: sentTime,
		Events:   []*pb.Event{},
	}
	mockCollector := new(collector.MockCollector)
	contextWithIDGroup := metadata.NewIncomingContext(ctx, meta)
	mockCollector.On("Collect", contextWithIDGroup, mock.Anything).Return(nil).Once()

	mockCollector.On("Collect", contextWithIDGroup, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		args.Get(1).(*collector.CollectRequest).AckFunc(nil)
	})

	mockCollector.On("Collect", contextWithIDGroup, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		args.Get(1).(*collector.CollectRequest).AckFunc(fmt.Errorf("simulated error"))
	})

	metaWithoutGroup := metadata.MD{}
	metaWithoutGroup.Set(config.Server.Websocket.Conn.IDHeader, "1235")
	contextWithoutGroup := metadata.NewIncomingContext(ctx, metaWithoutGroup)
	mockCollector.On("Collect", contextWithoutGroup, mock.Anything).Return(nil).Once()

	tests := []struct {
		name    string
		fields  fields
		args    args
		ack     config.AckType
		want    *pb.SendEventResponse
		wantErr bool
	}{
		{
			name: "Sending normal event",
			fields: fields{
				C:   mockCollector,
				ack: config.AckTypeAsync,
			},
			args: args{
				ctx: contextWithIDGroup,
				req: req,
			},
			want: &pb.SendEventResponse{
				Status:   pb.Status_STATUS_SUCCESS,
				Code:     pb.Code_CODE_OK,
				SentTime: sentTime.Seconds,
				Data: map[string]string{
					"req_guid": req.ReqGuid,
				},
			},
		},
		{
			name: "Sending normal event with synchronous ack",
			fields: fields{
				C:   mockCollector,
				ack: config.AckTypeSync,
			},
			args: args{
				ctx: contextWithIDGroup,
				req: req,
			},
			want: &pb.SendEventResponse{
				Status:   pb.Status_STATUS_SUCCESS,
				Code:     pb.Code_CODE_OK,
				SentTime: sentTime.Seconds,
				Data: map[string]string{
					"req_guid": req.ReqGuid,
				},
			},
		},
		{
			name: "Sending normal event with synchronous ack and collector error",
			fields: fields{
				C:   mockCollector,
				ack: config.AckTypeSync,
			},
			args: args{
				ctx: contextWithIDGroup,
				req: req,
			},
			want: &pb.SendEventResponse{
				Status:   pb.Status_STATUS_ERROR,
				Code:     pb.Code_CODE_INTERNAL_ERROR,
				SentTime: sentTime.Seconds,
				Data: map[string]string{
					"req_guid": req.ReqGuid,
				},
			},
		},
		{
			name: "Sending without group",
			fields: fields{
				C:   mockCollector,
				ack: config.AckTypeAsync,
			},
			args: args{
				ctx: contextWithoutGroup,
				req: req,
			},
			want: &pb.SendEventResponse{
				Status:   pb.Status_STATUS_SUCCESS,
				Code:     pb.Code_CODE_OK,
				SentTime: sentTime.Seconds,
				Data: map[string]string{
					"req_guid": req.ReqGuid,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				C:       tt.fields.C,
				ackType: tt.fields.ack,
			}
			got, err := h.SendEvent(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.SendEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.SendEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
