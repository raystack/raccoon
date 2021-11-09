package rest

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"raccoon/config"
	"raccoon/logger"
	"raccoon/metrics"
	"raccoon/pkg/collection"
	"raccoon/pkg/deserialization"
	"raccoon/pkg/identification"
	pb "raccoon/pkg/proto"
	"raccoon/pkg/serialization"
	"time"
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
	serDeMap map[string]*serDe
}

func NewHandler() *Handler {
	serDeMap := make(map[string]*serDe)
	serDeMap[ContentJSON] = &serDe{
		serializer:   serialization.JSONSerializer(),
		deserializer: deserialization.JSONDeserializer(),
	}

	serDeMap[ContentProto] = &serDe{
		serializer:   serialization.ProtoSerilizer(),
		deserializer: deserialization.ProtoDeserilizer(),
	}
	return &Handler{
		serDeMap: serDeMap,
	}
}

func (h *Handler) GetRESTAPIHandler(c collection.Collector) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		rw.Header().Set("Content-Type", contentType)

		res := &Response{}
		serde, ok := h.serDeMap[contentType]

		if !ok {
			metrics.Increment("batches_read_total", "status=failed,reason=unknowncontentype")
			logger.Errorf("[rest.GetRESTAPIHandler] invalid content type %s", contentType)
			rw.WriteHeader(http.StatusBadRequest)
			_, err := res.SetCode(pb.Code_BAD_REQUEST).SetStatus(pb.Status_ERROR).SetReason("no body present").
				SetSentTime(time.Now().Unix()).Write(rw, serialization.JSONSerializer())
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
			_, err := res.SetCode(pb.Code_BAD_REQUEST).SetStatus(pb.Status_ERROR).SetReason("no body present").
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
			_, err := res.SetCode(pb.Code_INTERNAL_ERROR).SetStatus(pb.Status_ERROR).SetReason("deserialization failure").
				SetSentTime(time.Now().Unix()).Write(rw, s)
			if err != nil {
				logger.Errorf("[restGetRESTAPIHandler] %s error sending error response: %v", identifier, err)
			}
			return
		}

		timeConsumed := time.Now()
		metrics.Count("events_rx_bytes_total", len(b), fmt.Sprintf("conn_group=%s", identifier.Group))
		req := &pb.EventRequest{}

		if err := d.Deserialize(b, req); err != nil {
			logger.Errorf("[rest.GetRESTAPIHandler] error while calling d.Deserialize() for %s, error: %s", identifier, err)
			metrics.Increment("batches_read_total", fmt.Sprintf("status=failed,reason=serde,conn_group=%s", identifier.Group))
			rw.WriteHeader(http.StatusBadRequest)
			_, err := res.SetCode(pb.Code_BAD_REQUEST).SetStatus(pb.Status_ERROR).SetReason("deserialization failure").
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
}
