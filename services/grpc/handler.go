package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/odpf/raccoon/collection"
	"github.com/odpf/raccoon/config"
	"github.com/odpf/raccoon/identification"
	"github.com/odpf/raccoon/logger"
	"github.com/odpf/raccoon/metrics"
	pb "github.com/odpf/raccoon/proto"
	"google.golang.org/grpc/metadata"
)

type Handler struct {
	C collection.Collector
	pb.UnimplementedEventServiceServer
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
	case 0:
		responseChannel <- &pb.SendEventResponse{
			Status:   pb.Status_STATUS_SUCCESS,
			Code:     pb.Code_CODE_OK,
			SentTime: time.Now().Unix(),
			Data: map[string]string{
				"req_guid": reqGuid,
			},
		}
		return nil
	case 1:
		return func(err error) {
			if err != nil {
				logger.Error(fmt.Sprintf("[grpc.Ack] publish message failed for %s: %v", connGroup, err))
				responseChannel <- &pb.SendEventResponse{
					Status:   pb.Status_STATUS_ERROR,
					Code:     pb.Code_CODE_UNSPECIFIED,
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
