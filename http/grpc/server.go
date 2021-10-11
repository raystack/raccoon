package grpc

import (
	"context"
	"fmt"
	"net/http"
	"raccoon/config"
	"raccoon/logger"
	"raccoon/metrics"
	"raccoon/pkg/collection"
	"raccoon/pkg/identification"
	pb "raccoon/pkg/proto"
	"time"

	"google.golang.org/grpc/metadata"
)

type Handler struct {
	pb.UnimplementedEventServiceServer
}

func (h *Handler) SendEvent(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	identifier := identification.Identifier{
		ID:    metadata.Get(config.ServerWs.ConnIDHeader)[0],
		Group: metadata.Get(config.ServerWs.ConnGroupHeader)[0],
	}

	timeConsumed := time.Now()
	metrics.Count("events_rx_bytes_total", len(), fmt.Sprintf("conn_group=%s", identifier.Group))
	req := &pb.EventRequest{}

	if err := d.Deserialize(b, req); err != nil {
		logger.Errorf("[rest.GetRESTAPIHandler] error while calling d.Deserialize() for %s, error: %s", identifier, err)
		metrics.Increment("batches_read_total", fmt.Sprintf("status=failed,reason=deserr,conn_group=%s", identifier.Group))
		rw.WriteHeader(http.StatusInternalServerError)
		_, err := res.SetCode(pb.Code_INTERNAL_ERROR).SetStatus(pb.Status_ERROR).SetReason("deserialization failure").
			SetSentTime(time.Now().Unix()).Write(rw, s)
		if err != nil {
			logger.Errorf("[restGetRESTAPIHandler] %s error sending error response: %v", identifier, err)
		}
		return
	}

	metrics.Increment("batches_read_total", fmt.Sprintf("status=success,conn_group=%s", identifier.Group))
	metrics.Count("events_rx_total", len(req.Events), fmt.Sprintf("conn_group=%s", identifier.Group))

	c.Collect(r.Context(), &collection.CollectRequest{
		ConnectionIdentifier: &identifier,
		TimeConsumed:         timeConsumed,
		EventRequest:         req,
	})

	_, err = res.SetCode(pb.Code_OK).SetStatus(pb.Status_SUCCESS).SetSentTime(time.Now().Unix()).
		SetDataMap(map[string]string{"req_guid": req.ReqGuid}).Write(rw, s)
	if err != nil {
		logger.Errorf("[restGetRESTAPIHandler] %s error sending error response: %v", identifier, err)
	}

}
