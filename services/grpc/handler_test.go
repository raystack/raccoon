package grpc

import (
	"context"
	"raccoon/collection"
	"raccoon/config"
	"raccoon/logger"
	"raccoon/metrics"
	pb "raccoon/proto"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type void struct{}

func (v void) Write(_ []byte) (int, error) {
	return 0, nil
}

func TestHandler_SendEvent(t *testing.T) {
	type fields struct {
		C                               collection.Collector
		UnimplementedEventServiceServer pb.UnimplementedEventServiceServer
	}
	type args struct {
		ctx context.Context
		req *pb.SendEventRequest
	}

	logger.SetOutput(void{})
	metrics.SetVoid()
	collector := new(collection.MockCollector)
	ctx := context.Background()
	meta := metadata.MD{}
	meta.Set(config.ServerWs.ConnGroupHeader, "group")
	meta.Set(config.ServerWs.ConnIDHeader, "1235")
	sentTime := timestamppb.Now()
	req := &pb.SendEventRequest{
		ReqGuid:  "abcd",
		SentTime: sentTime,
		Events:   []*pb.Event{},
	}
	contextWithIDGroup := metadata.NewIncomingContext(ctx, meta)
	collector.On("Collect", contextWithIDGroup, mock.Anything).Return(nil)

	metaWithoutGroup := metadata.MD{}
	metaWithoutGroup.Set(config.ServerWs.ConnIDHeader, "1235")
	contextWithoutGroup := metadata.NewIncomingContext(ctx, metaWithoutGroup)
	collector.On("Collect", contextWithoutGroup, mock.Anything).Return(nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.SendEventResponse
		wantErr bool
	}{
		{
			name: "Sending normal event",
			fields: fields{
				C: collector,
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
			name: "Sending without group",
			fields: fields{
				C: collector,
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
				C: tt.fields.C,
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
