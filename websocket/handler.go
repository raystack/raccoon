package websocket

import (
	"fmt"
	"net/http"
	"raccoon/logger"
	"raccoon/metrics"
	"raccoon/websocket/connection"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	pb "raccoon/websocket/proto"
)

type Handler struct {
	upgrader      *connection.Upgrader
	bufferChannel chan EventsBatch
	PingChannel   chan connection.Conn
}
type EventsBatch struct {
	ConnIdentifer connection.Identifer
	EventReq      *pb.EventRequest
	TimeConsumed  time.Time
	TimePushed    time.Time
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

//HandlerWSEvents handles the upgrade and the events sent by the peers
func (h *Handler) HandlerWSEvents(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r)
	if err != nil {
		logger.Debugf("[websocket.Handler] %v", err)
		return
	}
	defer conn.Close()
	h.PingChannel <- conn

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
				websocket.CloseNoStatusReceived,
				websocket.CloseAbnormalClosure) {
				logger.Error(fmt.Sprintf("[websocket.Handler] %s closed abruptly: %v", conn.Identifier, err))
				metrics.Increment("batches_read_total", "status=failed,reason=closeerror")
				break
			}

			metrics.Increment("batches_read_total", "status=failed,reason=unknown")
			logger.Error(fmt.Sprintf("[websocket.Handler] reading message failed. Unknown failure for %s: %v", conn.Identifier, err)) //no connection issue here
			break
		}
		timeConsumed := time.Now()
		metrics.Count("events_rx_bytes_total", len(message), "")
		payload := &pb.EventRequest{}
		err = proto.Unmarshal(message, payload)
		if err != nil {
			logger.Error(fmt.Sprintf("[websocket.Handler] reading message failed for %s: %v", conn.Identifier, err))
			metrics.Increment("batches_read_total", "status=failed,reason=serde")
			badrequest := createBadrequestResponse(err)
			conn.WriteMessage(websocket.BinaryMessage, badrequest)
			continue
		}
		metrics.Increment("batches_read_total", "status=success")
		metrics.Count("events_rx_total", len(payload.Events), "")

		h.bufferChannel <- EventsBatch{
			ConnIdentifer: conn.Identifier,
			EventReq:      payload,
			TimeConsumed:  timeConsumed,
			TimePushed:    (time.Now()),
		}

		resp := createSuccessResponse(payload.ReqGuid)
		success, _ := proto.Marshal(resp)
		conn.WriteMessage(websocket.BinaryMessage, success)
	}
}
