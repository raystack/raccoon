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
	serializer   serialization.SerializeFunc
	deserializer deserialization.DeserializeFunc
}
type Handler struct {
	serDeMap  map[string]*serDe
	collector collection.Collector
}

func NewHandler(collector collection.Collector) *Handler {
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
		metrics.Increment("batches_read_total", "status=failed,reason=unknowncontentype")
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

	if err := d(b, req); err != nil {
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

	resChannel := make(chan struct{}, 1)
	h.collector.Collect(r.Context(), &collection.CollectRequest{
		ConnectionIdentifier: identifier,
		TimeConsumed:         timeConsumed,
		SendEventRequest:     req,
		AckFunc:              h.Ack(rw, resChannel, s, req.ReqGuid, identifier.Group),
	})
	<-resChannel
}

func (h *Handler) Ack(rw http.ResponseWriter, resChannel chan struct{}, s serialization.SerializeFunc, reqGuid string, connGroup string) collection.AckFunc {
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
		metrics.Count("events_rx_bytes_total", len(e.EventBytes), fmt.Sprintf("conn_group=%s,event_type=%s", group, e.Type))
		metrics.Increment("events_rx_total", fmt.Sprintf("conn_group=%s,event_type=%s", group, e.Type))
	}
}
