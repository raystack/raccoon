package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	pbgrpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/raccoon/v1beta1/raccoonv1beta1grpc"
	pb "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/raccoon/v1beta1"
	"github.com/goto/raccoon/collection"
	"github.com/goto/raccoon/config"
	"github.com/goto/raccoon/identification"
	"github.com/goto/raccoon/logger"
	"github.com/goto/raccoon/metrics"
	"google.golang.org/grpc/metadata"
)

type Handler struct {
	C collection.Collector
	pbgrpc.UnimplementedEventServiceServer
}

func (h *Handler) SendEvent(ctx context.Context, req *pb.SendEventRequest) (*pb.SendEventResponse, error) {
	metadata, _ := metadata.FromIncomingContext(ctx)
	groups := metadata.Get(config.ServerWs.ConnGroupHeader)
	var group string
	if len(groups) > 0 {
		group = groups[0]
	} else {
		group = config.ServerWs.ConnGroupDefault
	}

	var id string
	ids := metadata.Get(config.ServerWs.ConnIDHeader)

	if len(ids) > 0 {
		id = ids[0]
	} else {
		return nil, errors.New("connection id header missing")
	}

	identifier := identification.Identifier{
		ID:    id,
		Group: group,
	}

	timeConsumed := time.Now()

	metrics.Increment("batches_read_total", fmt.Sprintf("status=success,conn_group=%s", identifier.Group))
	h.sendEventCounters(req.Events, identifier.Group)

	responseChannel := make(chan *pb.SendEventResponse, 1)
	h.C.Collect(ctx, &collection.CollectRequest{
		ConnectionIdentifier: identifier,
		TimeConsumed:         timeConsumed,
		SendEventRequest:     req,
		AckFunc:              h.Ack(responseChannel, req.ReqGuid, identifier.Group),
	})
	return <-responseChannel, nil

}

func (h *Handler) Ack(responseChannel chan *pb.SendEventResponse, reqGuid, connGroup string) collection.AckFunc {
	switch config.Event.Ack {
	case config.Asynchronous:
		responseChannel <- &pb.SendEventResponse{
			Status:   pb.Status_STATUS_SUCCESS,
			Code:     pb.Code_CODE_OK,
			SentTime: time.Now().Unix(),
			Data: map[string]string{
				"req_guid": reqGuid,
			},
		}
		return nil
	case config.Synchronous:
		return func(err error) {
			if err != nil {
				logger.Errorf("[grpc.Ack] publish message failed for %s: %v", connGroup, err)
				responseChannel <- &pb.SendEventResponse{
					Status:   pb.Status_STATUS_ERROR,
					Code:     pb.Code_CODE_INTERNAL_ERROR,
					SentTime: time.Now().Unix(),
					Data: map[string]string{
						"req_guid": reqGuid,
					},
				}
				return
			}
			responseChannel <- &pb.SendEventResponse{
				Status:   pb.Status_STATUS_SUCCESS,
				Code:     pb.Code_CODE_OK,
				SentTime: time.Now().Unix(),
				Data: map[string]string{
					"req_guid": reqGuid,
				},
			}
		}
	default:
		responseChannel <- &pb.SendEventResponse{
			Status:   pb.Status_STATUS_SUCCESS,
			Code:     pb.Code_CODE_OK,
			SentTime: time.Now().Unix(),
			Data: map[string]string{
				"req_guid": reqGuid,
			},
		}
		return nil
	}
}

func (h *Handler) sendEventCounters(events []*pb.Event, group string) {
	for _, e := range events {
		metrics.Increment("events_rx_total", fmt.Sprintf("conn_group=%s,event_type=%s", group, e.Type))
	}
}
