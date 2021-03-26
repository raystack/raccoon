package websocket

import (
	"fmt"
	"net/http"
	"raccoon/logger"
	"raccoon/metrics"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	pb "raccoon/websocket/proto"
)

type Handler struct {
	websocketUpgrader websocket.Upgrader
	bufferChannel     chan EventsBatch
	user              *User
	PongWaitInterval  time.Duration
	WriteWaitInterval time.Duration
	PingChannel       chan connection
	UniqConnIDHeader  string
}

type EventsBatch struct {
	UniqConnID   string
	EventReq     *pb.EventRequest
	TimeConsumed time.Time
	TimePushed   time.Time
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

//HandlerWSEvents handles the upgrade and the events sent by the peers
func (wsHandler *Handler) HandlerWSEvents(w http.ResponseWriter, r *http.Request) {
	uniqConnID := r.Header.Get(wsHandler.UniqConnIDHeader)
	connectedTime := time.Now()
	logger.Debug(fmt.Sprintf("UniqConnID %s connected at %v", uniqConnID, connectedTime))
	conn, err := wsHandler.websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("[websocket.Handler] Failed to upgrade connection  UniqConnID: %s : %v", uniqConnID, err))
		metrics.Increment("user_connection_failure_total", "reason=ugfailure")
		return
	}
	defer conn.Close()

	if wsHandler.user.Exists(uniqConnID) {
		logger.Errorf("[websocket.Handler] Disconnecting %v, already connected", uniqConnID)
		duplicateConnResp := createEmptyErrorResponse(pb.Code_MAX_USER_LIMIT_REACHED)

		conn.WriteMessage(websocket.BinaryMessage, duplicateConnResp)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "Duplicate connection"))
		metrics.Increment("user_connection_failure_total", "reason=exists")
		return
	}
	if wsHandler.user.HasReachedLimit() {
		logger.Errorf("[websocket.Handler] Disconnecting %v, max connection reached", uniqConnID)
		maxConnResp := createEmptyErrorResponse(pb.Code_MAX_CONNECTION_LIMIT_REACHED)
		conn.WriteMessage(websocket.BinaryMessage, maxConnResp)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "Max connection reached"))
		metrics.Increment("user_connection_failure_total", "reason=serverlimit")
		return
	}
	wsHandler.user.Store(uniqConnID)
	defer wsHandler.user.Remove(uniqConnID)
	defer calculateSessionTime(uniqConnID, connectedTime)

	setUpControlHandlers(conn, uniqConnID, wsHandler.PongWaitInterval, wsHandler.WriteWaitInterval)
	wsHandler.PingChannel <- connection{
		uniqConnID: uniqConnID,
		conn:       conn,
	}
	metrics.Increment("user_connection_success_total", "")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
				websocket.CloseNoStatusReceived,
				websocket.CloseAbnormalClosure) {
				logger.Error(fmt.Sprintf("[websocket.Handler] UniqConnID %s Connection Closed Abruptly: %v", uniqConnID, err))
				metrics.Increment("batches_read_total", "status=failed,reason=closeerror")
				break
			}

			metrics.Increment("batches_read_total", "status=failed,reason=unknown")
			logger.Error(fmt.Sprintf("[websocket.Handler] Reading message failed. Unknown failure: %v  User ID: %s ", err, uniqConnID)) //no connection issue here
			break
		}
		timeConsumed := time.Now()
		metrics.Count("events_rx_bytes_total", len(message), "")
		payload := &pb.EventRequest{}
		err = proto.Unmarshal(message, payload)
		if err != nil {
			logger.Error(fmt.Sprintf("[websocket.Handler] Reading message failed. %v  UniqConnID: %s ", err, uniqConnID))
			metrics.Increment("batches_read_total", "status=failed,reason=serde")
			badrequest := createBadrequestResponse(err)
			conn.WriteMessage(websocket.BinaryMessage, badrequest)
			continue
		}
		metrics.Increment("batches_read_total", "status=success")
		metrics.Count("events_rx_total", len(payload.Events), "")

		wsHandler.bufferChannel <- EventsBatch{
			UniqConnID:   uniqConnID,
			EventReq:     payload,
			TimeConsumed: timeConsumed,
			TimePushed:   (time.Now()),
		}

		resp := createSuccessResponse(payload.ReqGuid)
		success, _ := proto.Marshal(resp)
		conn.WriteMessage(websocket.BinaryMessage, success)
	}
}

func calculateSessionTime(uniqConnID string, connectedAt time.Time) {
	connectionTime := time.Now().Sub(connectedAt)
	logger.Debug(fmt.Sprintf("[websocket.calculateSessionTime] UniqConnID: %s, total time connected in minutes: %v", uniqConnID, connectionTime.Minutes()))
	metrics.Timing("user_session_duration_milliseconds", connectionTime.Milliseconds(), "")
}

func setUpControlHandlers(conn *websocket.Conn, uniqConnID string,
	PongWaitInterval time.Duration, WriteWaitInterval time.Duration) {
	//expects the client to send a ping, mark this channel as idle timed out post the deadline
	conn.SetReadDeadline(time.Now().Add(PongWaitInterval))
	conn.SetPongHandler(func(string) error {
		// extends the read deadline since we have received this pong on this channel
		conn.SetReadDeadline(time.Now().Add(PongWaitInterval))
		return nil
	})

	conn.SetPingHandler(func(s string) error {
		logger.Debug(fmt.Sprintf("Client connection with UniqConnID: %s Pinged", uniqConnID))
		if err := conn.WriteControl(websocket.PongMessage, []byte(s), time.Now().Add(WriteWaitInterval)); err != nil {
			metrics.Increment("server_pong_failure_total", "")
			logger.Debug(fmt.Sprintf("Failed to send pong event: %s UniqConnID: %s", err, uniqConnID))
		}
		return nil
	})
}
