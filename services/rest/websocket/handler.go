package websocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/raystack/raccoon/collector"
	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/deserialization"
	"github.com/raystack/raccoon/logger"
	"github.com/raystack/raccoon/metrics"
	pb "github.com/raystack/raccoon/proto"
	"github.com/raystack/raccoon/serialization"
	"github.com/raystack/raccoon/services/rest/websocket/connection"
)

type serDe struct {
	serializer   serialization.SerializeFunc
	deserializer deserialization.DeserializeFunc
}
type Handler struct {
	upgrader    *connection.Upgrader
	serdeMap    map[int]*serDe
	collector   collector.Collector
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

func NewHandler(pingC chan connection.Conn, collector collector.Collector) *Handler {
	ugConfig := connection.UpgraderConfig{
		ReadBufferSize:    config.Server.Websocket.ReadBufferSize,
		WriteBufferSize:   config.Server.Websocket.WriteBufferSize,
		CheckOrigin:       config.Server.Websocket.CheckOrigin,
		MaxUser:           config.Server.Websocket.ServerMaxConn,
		PongWaitInterval:  time.Duration(config.Server.Websocket.PongWaitIntervalMS) * time.Millisecond,
		WriteWaitInterval: time.Duration(config.Server.Websocket.WriteWaitIntervalMS) * time.Millisecond,
		ConnIDHeader:      config.Server.Websocket.Conn.IDHeader,
		ConnGroupHeader:   config.Server.Websocket.Conn.GroupHeader,
		ConnGroupDefault:  config.Server.Websocket.Conn.GroupDefault,
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
				metrics.Increment("batches_read_total", map[string]string{"status": "failed", "reason": "closeerror", "conn_group": conn.Identifier.Group})
				break
			}
			metrics.Increment("batches_read_total", map[string]string{"status": "failed", "reason": "unknown", "conn_group": conn.Identifier.Group})
			logger.Error(fmt.Sprintf("[websocket.Handler] reading message failed. Unknown failure for %s: %v", conn.Identifier, err)) //no connection issue here
			break
		}

		timeConsumed := time.Now()
		payload := &pb.SendEventRequest{}
		serde := h.serdeMap[messageType]
		d, s := serde.deserializer, serde.serializer
		if err := d(message, payload); err != nil {
			logger.Error(fmt.Sprintf("[websocket.Handler] reading message failed for %s: %v", conn.Identifier, err))
			metrics.Increment("batches_read_total", map[string]string{"status": "failed", "reason": "serde", "conn_group": conn.Identifier.Group})
			writeBadRequestResponse(conn, s, messageType, payload.ReqGuid, err)
			continue
		}
		if config.Server.Batch.DedupEnabled {
			// avoiding processing the same active connection's duplicate events.
			if h.upgrader.Table.HasBatch(conn.Identifier, payload.ReqGuid) {
				metrics.Increment("events_duplicate_total", map[string]string{"reason": "duplicate", "conn_group": conn.Identifier.Group})
				writeSuccessResponse(conn, s, messageType, payload.ReqGuid)
				continue
			}
			h.upgrader.Table.StoreBatch(conn.Identifier, payload.ReqGuid)
		}

		metrics.Increment("batches_read_total", map[string]string{"status": "success", "conn_group": conn.Identifier.Group, "reason": "NA"})
		h.sendEventCounters(payload.Events, conn.Identifier.Group)

		h.collector.Collect(r.Context(), &collector.CollectRequest{
			ConnectionIdentifier: conn.Identifier,
			TimeConsumed:         timeConsumed,
			SendEventRequest:     payload,
			AckFunc:              h.Ack(conn, AckChan, s, messageType, payload.ReqGuid, timeConsumed),
		})
	}
}

func (h *Handler) Ack(conn connection.Conn, resChannel chan AckInfo, s serialization.SerializeFunc, messageType int, reqGuid string, timeConsumed time.Time) collector.AckFunc {
	switch config.Event.Ack {
	case config.Asynchronous:
		writeSuccessResponse(conn, s, messageType, reqGuid)
		return nil
	case config.Synchronous:
		return func(err error) {
			if config.Server.Batch.DedupEnabled {
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
		metrics.Count("events_rx_bytes_total", int64(len(e.EventBytes)), map[string]string{"conn_group": group, "event_type": e.Type})
		metrics.Increment("events_rx_total", map[string]string{"conn_group": group, "event_type": e.Type})
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
