package websocket

import (
	"fmt"
	"net/http"
	"time"

	pb "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/raccoon/v1beta1"
	"github.com/gorilla/websocket"
	"github.com/goto/raccoon/collection"
	"github.com/goto/raccoon/config"
	"github.com/goto/raccoon/deserialization"
	"github.com/goto/raccoon/logger"
	"github.com/goto/raccoon/metrics"
	"github.com/goto/raccoon/serialization"
	"github.com/goto/raccoon/services/rest/websocket/connection"
)

type serDe struct {
	serializer   serialization.SerializeFunc
	deserializer deserialization.DeserializeFunc
}
type Handler struct {
	upgrader    *connection.Upgrader
	serdeMap    map[int]*serDe
	collector   collection.Collector
	PingChannel chan connection.Conn
}

func getSerDeMap() map[int]*serDe {
	serDeMap := make(map[int]*serDe)
	serDeMap[websocket.BinaryMessage] = &serDe{
		serializer:   serialization.SerializeProto,
		deserializer: deserialization.DeserializeProto,
	}

	serDeMap[websocket.TextMessage] = &serDe{
		serializer:   serialization.SerializeJSON,
		deserializer: deserialization.DeserializeJSON,
	}
	return serDeMap
}

func NewHandler(pingC chan connection.Conn, collector collection.Collector) *Handler {
	ugConfig := connection.UpgraderConfig{
		ReadBufferSize:    config.ServerWs.ReadBufferSize,
		WriteBufferSize:   config.ServerWs.WriteBufferSize,
		CheckOrigin:       config.ServerWs.CheckOrigin,
		MaxUser:           config.ServerWs.ServerMaxConn,
		PongWaitInterval:  config.ServerWs.PongWaitInterval,
		WriteWaitInterval: config.ServerWs.WriteWaitInterval,
		ConnIDHeader:      config.ServerWs.ConnIDHeader,
		ConnGroupHeader:   config.ServerWs.ConnGroupHeader,
		ConnGroupDefault:  config.ServerWs.ConnGroupDefault,
	}

	upgrader := connection.NewUpgrader(ugConfig)
	return &Handler{
		upgrader:    upgrader,
		serdeMap:    getSerDeMap(),
		PingChannel: pingC,
		collector:   collector,
	}
}

func (h *Handler) Table() *connection.Table {
	return h.upgrader.Table
}

// HandlerWSEvents handles the upgrade and the events sent by the peers
func (h *Handler) HandlerWSEvents(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r)
	if err != nil {
		logger.Errorf("[websocket.Handler] %v", err)
		return
	}
	defer conn.Close()
	h.PingChannel <- conn
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
				websocket.CloseNoStatusReceived,
				websocket.CloseAbnormalClosure) {
				logger.Error(fmt.Sprintf("[websocket.Handler] %s closed abruptly: %v", conn.Identifier, err))
				metrics.Increment("batches_read_total", fmt.Sprintf("status=failed,reason=closeerror,conn_group=%s", conn.Identifier.Group))
				break
			}
			metrics.Increment("batches_read_total", fmt.Sprintf("status=failed,reason=unknown,conn_group=%s", conn.Identifier.Group))
			logger.Error(fmt.Sprintf("[websocket.Handler] reading message failed. Unknown failure for %s: %v", conn.Identifier, err)) //no connection issue here
			break
		}

		timeConsumed := time.Now()
		payload := &pb.SendEventRequest{}
		serde := h.serdeMap[messageType]
		d, s := serde.deserializer, serde.serializer
		if err := d(message, payload); err != nil {
			logger.Error(fmt.Sprintf("[websocket.Handler] reading message failed for %s: %v", conn.Identifier, err))
			metrics.Increment("batches_read_total", fmt.Sprintf("status=failed,reason=serde,conn_group=%s", conn.Identifier.Group))
			writeBadRequestResponse(conn, s, messageType, payload.ReqGuid, err)
			continue
		}
		if config.Server.DedupEnabled {
			// avoiding processing the same active connection's duplicate events.
			if h.upgrader.Table.HasBatch(conn.Identifier, payload.ReqGuid) {
				metrics.Increment("events_duplicate_total", fmt.Sprintf("reason=duplicate,conn_group=%s", conn.Identifier.Group))
				writeSuccessResponse(conn, s, messageType, payload.ReqGuid)
				continue
			}
			h.upgrader.Table.StoreBatch(conn.Identifier, payload.ReqGuid)
		}

		metrics.Increment("batches_read_total", fmt.Sprintf("status=success,conn_group=%s", conn.Identifier.Group))
		h.sendEventCounters(payload.Events, conn.Identifier.Group)

		h.collector.Collect(r.Context(), &collection.CollectRequest{
			ConnectionIdentifier: conn.Identifier,
			TimeConsumed:         timeConsumed,
			SendEventRequest:     payload,
			AckFunc:              h.Ack(conn, AckChan, s, messageType, payload.ReqGuid, timeConsumed),
		})
	}
}

func (h *Handler) Ack(conn connection.Conn, resChannel chan AckInfo, s serialization.SerializeFunc, messageType int, reqGuid string, timeConsumed time.Time) collection.AckFunc {
	switch config.Event.Ack {
	case config.Asynchronous:
		writeSuccessResponse(conn, s, messageType, reqGuid)
		return nil
	case config.Synchronous:
		return func(err error) {
			if config.Server.DedupEnabled {
				if err != nil {
					h.upgrader.Table.RemoveBatch(conn.Identifier, reqGuid)
				}
			}
			AckChan <- AckInfo{
				MessageType:     messageType,
				RequestGuid:     reqGuid,
				Err:             err,
				Conn:            conn,
				serializer:      h.serdeMap[messageType].serializer,
				TimeConsumed:    timeConsumed,
				AckTimeConsumed: time.Now(),
			}
		}
	default:
		writeSuccessResponse(conn, s, messageType, reqGuid)
		return nil
	}
}

func (h *Handler) sendEventCounters(events []*pb.Event, group string) {
	for _, e := range events {
		metrics.Count("events_rx_bytes_total", len(e.EventBytes), fmt.Sprintf("conn_group=%s,event_type=%s", group, e.Type))
		metrics.Increment("events_rx_total", fmt.Sprintf("conn_group=%s,event_type=%s", group, e.Type))
	}
}

func writeSuccessResponse(conn connection.Conn, serialize serialization.SerializeFunc, messageType int, requestGUID string) {
	response := &pb.SendEventResponse{
		Status:   pb.Status_STATUS_SUCCESS,
		Code:     pb.Code_CODE_OK,
		SentTime: time.Now().Unix(),
		Reason:   "",
		Data: map[string]string{
			"req_guid": requestGUID,
		},
	}
	success, _ := serialize(response)
	conn.WriteMessage(messageType, success)
}

func writeBadRequestResponse(conn connection.Conn, serialize serialization.SerializeFunc, messageType int, reqGuid string, err error) {
	response := &pb.SendEventResponse{
		Status:   pb.Status_STATUS_ERROR,
		Code:     pb.Code_CODE_BAD_REQUEST,
		SentTime: time.Now().Unix(),
		Reason:   fmt.Sprintf("cannot deserialize request: %s", err),
		Data: map[string]string{
			"req_guid": reqGuid,
		},
	}

	failure, _ := serialize(response)
	conn.WriteMessage(messageType, failure)
}

func writeFailedResponse(conn connection.Conn, serialize serialization.SerializeFunc, messageType int, reqGuid string, err error) {
	response := &pb.SendEventResponse{
		Status:   pb.Status_STATUS_ERROR,
		Code:     pb.Code_CODE_INTERNAL_ERROR,
		SentTime: time.Now().Unix(),
		Reason:   fmt.Sprintf("cannot publish events: %s", err),
		Data: map[string]string{
			"req_guid": reqGuid,
		},
	}
	failure, _ := serialize(response)
	conn.WriteMessage(messageType, failure)
}
