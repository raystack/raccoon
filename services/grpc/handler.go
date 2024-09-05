package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/raystack/raccoon/collector"
	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/identification"
	"github.com/raystack/raccoon/logger"
	"github.com/raystack/raccoon/metrics"
	pb "github.com/raystack/raccoon/proto"
	"google.golang.org/grpc/metadata"
)

type Handler struct {
	C collector.Collector
	pb.UnimplementedEventServiceServer
}

func (h *Handler) SendEvent(ctx context.Context, req *pb.SendEventRequest) (*pb.SendEventResponse, error) {
	metadata, _ := metadata.FromIncomingContext(ctx)
	groups := metadata.Get(config.Server.Websocket.Conn.GroupHeader)
	var group string
	if len(groups) > 0 {
		group = groups[0]
	} else {
		group = config.Server.Websocket.Conn.GroupDefault
	}

	var id string
	ids := metadata.Get(config.Server.Websocket.Conn.IDHeader)

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

	metrics.Increment("batches_read_total", map[string]string{"status": "success", "conn_group": identifier.Group, "reason": "NA"})
	h.sendEventCounters(req.Events, identifier.Group)

	responseChannel := make(chan *pb.SendEventResponse, 1)
	h.C.Collect(ctx, &collector.CollectRequest{
		ConnectionIdentifier: identifier,
		TimeConsumed:         timeConsumed,
		SendEventRequest:     req,
		AckFunc:              h.Ack(responseChannel, req.ReqGuid, identifier.Group),
	})
	return <-responseChannel, nil
}

func (h *Handler) Ack(responseChannel chan *pb.SendEventResponse, reqGuid, connGroup string) collector.AckFunc {
	switch config.Event.Ack {
	case config.AckTypeAsync:
		responseChannel <- &pb.SendEventResponse{
			Status:   pb.Status_STATUS_SUCCESS,
			Code:     pb.Code_CODE_OK,
			SentTime: time.Now().Unix(),
			Data: map[string]string{
				"req_guid": reqGuid,
			},
		}
		return nil
	case config.AckTypeSync:
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
		metrics.Increment("events_rx_total", map[string]string{"conn_group": group, "event_type": e.Type})
	}
}
