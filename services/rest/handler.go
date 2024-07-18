package rest

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/raystack/raccoon/collector"
	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/deserialization"
	"github.com/raystack/raccoon/identification"
	"github.com/raystack/raccoon/logger"
	"github.com/raystack/raccoon/metrics"
	pb "github.com/raystack/raccoon/proto"
	"github.com/raystack/raccoon/serialization"
)

const (
	ContentJSON  = "application/json"
	ContentProto = "application/proto"
)

type serDe struct {
	serializer   serialization.SerializeFunc
	deserializer deserialization.DeserializeFunc
}
type Handler struct {
	serDeMap  map[string]*serDe
	collector collector.Collector
}

func NewHandler(collector collector.Collector) *Handler {
	serDeMap := make(map[string]*serDe)
	serDeMap[ContentJSON] = &serDe{
		serializer:   serialization.SerializeJSON,
		deserializer: deserialization.DeserializeJSON,
	}

	serDeMap[ContentProto] = &serDe{
		serializer:   serialization.SerializeProto,
		deserializer: deserialization.DeserializeProto,
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
		metrics.Increment("batches_read_total", map[string]string{"status": "failed", "reason": "unknowncontentype", "conn_group": "NA"})
		logger.Errorf("[rest.GetRESTAPIHandler] invalid content type %s", contentType)
		rw.WriteHeader(http.StatusBadRequest)
		_, err := res.SetCode(pb.Code_CODE_BAD_REQUEST).SetStatus(pb.Status_STATUS_ERROR).SetReason("invalid content type").
			SetSentTime(time.Now().Unix()).Write(rw, serialization.SerializeJSON)
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
		metrics.Increment("batches_read_total", map[string]string{"status": "failed", "reason": "emptybody", "conn_group": identifier.Group})
		logger.Errorf("[rest.GetRESTAPIHandler] %s no body", identifier)
		rw.WriteHeader(http.StatusBadRequest)
		_, err := res.SetCode(pb.Code_CODE_BAD_REQUEST).SetStatus(pb.Status_STATUS_ERROR).SetReason("no body present").
			SetSentTime(time.Now().Unix()).Write(rw, s)
		if err != nil {
			logger.Errorf("[rest.GetRESTAPIHandler] %s error sending response: %v", identifier, err)
		}
		return
	}

	defer io.Copy(io.Discard, r.Body)
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf(fmt.Sprintf("[rest.GetRESTAPIHandler] %s error reading request body, error: %v", identifier, err))
		metrics.Increment("batches_read_total", map[string]string{"status": "failed", "reason": "readerr", "conn_group": identifier.Group})
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

	if err := d(b, req); err != nil {
		logger.Errorf("[rest.GetRESTAPIHandler] error while calling d.Deserialize() for %s, error: %s", identifier, err)
		metrics.Increment("batches_read_total", map[string]string{"status": "failed", "reason": "serde", "conn_group": identifier.Group})
		rw.WriteHeader(http.StatusBadRequest)
		_, err := res.SetCode(pb.Code_CODE_BAD_REQUEST).SetStatus(pb.Status_STATUS_ERROR).SetReason("deserialization failure").
			SetSentTime(time.Now().Unix()).Write(rw, s)
		if err != nil {
			logger.Errorf("[restGetRESTAPIHandler] %s error sending error response: %v", identifier, err)
		}
		return
	}

	metrics.Increment("batches_read_total", map[string]string{"status": "success", "conn_group": identifier.Group, "reason": "NA"})
	h.sendEventCounters(req.Events, identifier.Group)

	resChannel := make(chan struct{}, 1)
	h.collector.Collect(r.Context(), &collector.CollectRequest{
		ConnectionIdentifier: identifier,
		TimeConsumed:         timeConsumed,
		SendEventRequest:     req,
		AckFunc:              h.Ack(rw, resChannel, s, req.ReqGuid, identifier.Group),
	})
	<-resChannel
}

func (h *Handler) Ack(rw http.ResponseWriter, resChannel chan struct{}, s serialization.SerializeFunc, reqGuid string, connGroup string) collector.AckFunc {
	res := &Response{
		SendEventResponse: &pb.SendEventResponse{},
	}
	switch config.Event.Ack {
	case config.Asynchronous:

		rw.WriteHeader(http.StatusOK)
		_, err := res.SetCode(pb.Code_CODE_OK).SetStatus(pb.Status_STATUS_SUCCESS).SetSentTime(time.Now().Unix()).
			SetDataMap(map[string]string{"req_guid": reqGuid}).Write(rw, s)
		if err != nil {
			logger.Errorf("[RESTAPIHandler.Ack] %s error sending error response: %v", connGroup, err)
		}
		resChannel <- struct{}{}
		return nil
	case config.Synchronous:
		return func(err error) {
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				logger.Errorf("[RESTAPIHandler.Ack] publish message failed for %s: %v", connGroup, err)
				_, err := res.SetCode(pb.Code_CODE_INTERNAL_ERROR).SetStatus(pb.Status_STATUS_ERROR).SetReason(fmt.Sprintf("cannot publish events: %s", err)).
					SetSentTime(time.Now().Unix()).SetDataMap(map[string]string{"req_guid": reqGuid}).Write(rw, s)
				if err != nil {
					logger.Errorf("[RESTAPIHandler] %s error sending error response: %v", connGroup, err)
				}
				resChannel <- struct{}{}
				return
			}
			rw.WriteHeader(http.StatusOK)
			_, err = res.SetCode(pb.Code_CODE_OK).SetStatus(pb.Status_STATUS_SUCCESS).SetSentTime(time.Now().Unix()).
				SetDataMap(map[string]string{"req_guid": reqGuid}).Write(rw, s)
			if err != nil {
				logger.Errorf("[RESTAPIHandler] %s error sending error response: %v", connGroup, err)
			}
			resChannel <- struct{}{}
		}
	default:
		rw.WriteHeader(http.StatusOK)
		_, err := res.SetCode(pb.Code_CODE_OK).SetStatus(pb.Status_STATUS_SUCCESS).SetSentTime(time.Now().Unix()).
			SetDataMap(map[string]string{"req_guid": reqGuid}).Write(rw, s)
		if err != nil {
			logger.Errorf("[RESTAPIHandler.Ack] %s error sending error response: %v", connGroup, err)
		}
		resChannel <- struct{}{}
		return nil
	}
}

func (h *Handler) sendEventCounters(events []*pb.Event, group string) {
	for _, e := range events {
		metrics.Count("events_rx_bytes_total", int64(len(e.EventBytes)), map[string]string{"conn_group": group, "event_type": e.Type})
		metrics.Increment("events_rx_total", map[string]string{"conn_group": group, "event_type": e.Type})
	}
}
