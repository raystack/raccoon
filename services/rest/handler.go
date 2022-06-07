package rest

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/odpf/raccoon/collection"
	"github.com/odpf/raccoon/config"
	"github.com/odpf/raccoon/deserialization"
	"github.com/odpf/raccoon/identification"
	"github.com/odpf/raccoon/logger"
	"github.com/odpf/raccoon/metrics"
	pb "github.com/odpf/raccoon/proto"
	"github.com/odpf/raccoon/serialization"
)

const (
	ContentJSON  = "application/json"
	ContentProto = "application/proto"
)

type serDe struct {
	serializer   serialization.Serializer
	deserializer deserialization.Deserializer
}
type Handler struct {
	serDeMap  map[string]*serDe
	collector collection.Collector
}

func NewHandler(collector collection.Collector) *Handler {
	serDeMap := make(map[string]*serDe)
	serDeMap[ContentJSON] = &serDe{
		serializer:   &serialization.JSONSerializer{},
		deserializer: &deserialization.JSONDeserializer{},
	}

	serDeMap[ContentProto] = &serDe{
		serializer:   &serialization.ProtoSerilizer{},
		deserializer: &deserialization.ProtoDeserilizer{},
	}
	return &Handler{
		serDeMap:  serDeMap,
		collector: collector,
	}
}

func (h *Handler) RESTAPIHandler(rw http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	rw.Header().Set("Content-Type", contentType)

	res := &Response{
		SendEventResponse: &pb.SendEventResponse{},
	}

	serde, ok := h.serDeMap[contentType]

	if !ok {
		metrics.Increment("batches_read_total", "status=failed,reason=unknowncontentype")
		logger.Errorf("[rest.GetRESTAPIHandler] invalid content type %s", contentType)
		rw.WriteHeader(http.StatusBadRequest)
		_, err := res.SetCode(pb.Code_CODE_BAD_REQUEST).SetStatus(pb.Status_STATUS_ERROR).SetReason("invalid content type").
			SetSentTime(time.Now().Unix()).Write(rw, &serialization.JSONSerializer{})
		if err != nil {
			logger.Errorf("[rest.GetRESTAPIHandler] error sending response: %v", err)
		}
		return
	}
	d, s := serde.deserializer, serde.serializer

	var group string
	group = r.Header.Get(config.ServerWs.ConnGroupHeader)
	if group == "" {
		group = config.ServerWs.ConnGroupDefault
	}
	identifier := identification.Identifier{
		ID:    r.Header.Get(config.ServerWs.ConnIDHeader),
		Group: group,
	}

	if r.Body == nil {
		metrics.Increment("batches_read_total", fmt.Sprintf("status=failed,reason=emptybody,conn_group=%s", identifier.Group))
		logger.Errorf("[rest.GetRESTAPIHandler] %s no body", identifier)
		rw.WriteHeader(http.StatusBadRequest)
		_, err := res.SetCode(pb.Code_CODE_BAD_REQUEST).SetStatus(pb.Status_STATUS_ERROR).SetReason("no body present").
			SetSentTime(time.Now().Unix()).Write(rw, s)
		if err != nil {
			logger.Errorf("[rest.GetRESTAPIHandler] %s error sending response: %v", identifier, err)
		}
		return
	}

	defer io.Copy(ioutil.Discard, r.Body)
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf(fmt.Sprintf("[rest.GetRESTAPIHandler] %s error reading request body, error: %v", identifier, err))
		metrics.Increment("batches_read_total", fmt.Sprintf("status=failed,reason=readerr,conn_group=%s", identifier.Group))
		rw.WriteHeader(http.StatusInternalServerError)
		_, err := res.SetCode(pb.Code_CODE_INTERNAL_ERROR).SetStatus(pb.Status_STATUS_ERROR).SetReason("deserialization failure").
			SetSentTime(time.Now().Unix()).Write(rw, s)
		if err != nil {
			logger.Errorf("[restGetRESTAPIHandler] %s error sending error response: %v", identifier, err)
		}
		return
	}

	timeConsumed := time.Now()
	req := &pb.SendEventRequest{}

	if err := d.Deserialize(b, req); err != nil {
		logger.Errorf("[rest.GetRESTAPIHandler] error while calling d.Deserialize() for %s, error: %s", identifier, err)
		metrics.Increment("batches_read_total", fmt.Sprintf("status=failed,reason=serde,conn_group=%s", identifier.Group))
		rw.WriteHeader(http.StatusBadRequest)
		_, err := res.SetCode(pb.Code_CODE_BAD_REQUEST).SetStatus(pb.Status_STATUS_ERROR).SetReason("deserialization failure").
			SetSentTime(time.Now().Unix()).Write(rw, s)
		if err != nil {
			logger.Errorf("[restGetRESTAPIHandler] %s error sending error response: %v", identifier, err)
		}
		return
	}

	metrics.Increment("batches_read_total", fmt.Sprintf("status=success,conn_group=%s", identifier.Group))
	h.sendEventCounters(req.Events, identifier.Group)

	h.collector.Collect(r.Context(), &collection.CollectRequest{
		ConnectionIdentifier: identifier,
		TimeConsumed:         timeConsumed,
		SendEventRequest:     req,
	})

	_, err = res.SetCode(pb.Code_CODE_OK).SetStatus(pb.Status_STATUS_SUCCESS).SetSentTime(time.Now().Unix()).
		SetDataMap(map[string]string{"req_guid": req.ReqGuid}).Write(rw, s)
	if err != nil {
		logger.Errorf("[restGetRESTAPIHandler] %s error sending error response: %v", identifier, err)
	}
}

func (h *Handler) sendEventCounters(events []*pb.Event, group string) {
	for _, e := range events {
		metrics.Count("events_rx_bytes_total", len(e.EventBytes), fmt.Sprintf("conn_group=%s,event_type=%s", group, e.Type))
		metrics.Increment("events_rx_total", fmt.Sprintf("conn_group=%s,event_type=%s", group, e.Type))
	}
}
