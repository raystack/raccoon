package grpc

import (
	"context"
	"fmt"
	"raccoon/config"
	"raccoon/metrics"
	"raccoon/pkg/collection"
	"raccoon/pkg/identification"
	pb "raccoon/pkg/proto"
	"time"

	"google.golang.org/grpc/metadata"
)

type Handler struct {
	C collection.Collector
	pb.UnimplementedEventServiceServer
}

func (h *Handler) SendEvent(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {
	metadata, _ := metadata.FromIncomingContext(ctx)
	groups := metadata.Get(config.ServerWs.ConnGroupHeader)
	var group string
	if len(groups) > 0 {
		group = groups[0]
	} else {
		group = config.ServerWs.ConnGroupDefault
	}
	identifier := identification.Identifier{
		ID:    metadata.Get(config.ServerWs.ConnIDHeader)[0],
		Group: group,
	}

	timeConsumed := time.Now()

	metrics.Increment("batches_read_total", fmt.Sprintf("status=success,conn_group=%s", identifier.Group))
	metrics.Count("events_rx_total", len(req.Events), fmt.Sprintf("conn_group=%s", identifier.Group))

	h.C.Collect(ctx, &collection.CollectRequest{
		ConnectionIdentifier: &identifier,
		TimeConsumed:         timeConsumed,
		EventRequest:         req,
	})

	return &pb.EventResponse{
		Status:   pb.Status_SUCCESS,
		Code:     pb.Code_OK,
		SentTime: time.Now().Unix(),
		Data: map[string]string{
			"req_guid": req.GetReqGuid(),
		},
	}, nil

}
