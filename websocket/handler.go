package websocket

import (
	"fmt"
	"net/http"
	"raccoon/logger"
	"raccoon/metrics"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"source.golabs.io/mobile/clickstream-go-proto/gojek/clickstream/de"
)

type Handler struct {
	websocketUpgrader websocket.Upgrader
	bufferChannel     chan EventsBatch
	user              *User
	PongWaitInterval  time.Duration
	WriteWaitInterval time.Duration
	PingChannel       chan connection
	UserIDHeader      string
}

type EventsBatch struct {
	UserID       string
	EventReq     *de.EventRequest
	TimeConsumed time.Time
	TimePushed   time.Time
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

//HandlerWSEvents handles the upgrade and the events sent by the peers
func (wsHandler *Handler) HandlerWSEvents(w http.ResponseWriter, r *http.Request) {
	UserID := r.Header.Get(wsHandler.UserIDHeader)
	connectedTime := time.Now()
	logger.Debug(fmt.Sprintf("UserID %s connected at %v", UserID, connectedTime))
	conn, err := wsHandler.websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("[websocket.Handler] Failed to upgrade connection  User ID: %s : %v", UserID, err))
		metrics.Increment("users.disconnected", "reason=ugfailure")
		return
	}
	defer conn.Close()

	if wsHandler.user.Exists(UserID) {
		logger.Errorf("[websocket.Handler] Disconnecting %v, already connected", UserID)
		duplicateConnResp := createEmptyErrorResponse(de.Code_MAX_USER_LIMIT_REACHED)

		conn.WriteMessage(websocket.BinaryMessage, duplicateConnResp)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "Duplicate connection"))
		metrics.Count("user.disconnected", 1, "reason=exists")
		return
	}
	if wsHandler.user.HasReachedLimit() {
		logger.Errorf("[websocket.Handler] Disconnecting %v, max connection reached", UserID)
		maxConnResp := createEmptyErrorResponse(de.Code_MAX_CONNECTION_LIMIT_REACHED)
		conn.WriteMessage(websocket.BinaryMessage, maxConnResp)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "Max connection reached"))
		metrics.Count("user.disconnected", 1, "reason=serverlimit")
		return
	}
	wsHandler.user.Store(UserID)
	defer wsHandler.user.Remove(UserID)
	defer calculateSessionTime(UserID, connectedTime)

	setUpControlHandlers(conn, UserID, wsHandler.PongWaitInterval, wsHandler.WriteWaitInterval)
	wsHandler.PingChannel <- connection{
		userID: UserID,
		conn:   conn,
	}
	metrics.Increment("user.connected", "")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
				websocket.CloseNoStatusReceived,
				websocket.CloseAbnormalClosure) {
				logger.Error(fmt.Sprintf("[websocket.Handler] UserID %s Connection Closed Abruptly: %v", UserID, err))
				metrics.Count("batches.read", 1, "status=failed,reason=closeerror")
				break
			}

			metrics.Count("batches.read", 1, "status=failed,reason=unknown")
			logger.Error(fmt.Sprintf("[websocket.Handler] Reading message failed. Unknown failure: %v  User ID: %s ", err, UserID)) //no connection issue here
			break
		}
		timeConsumed := time.Now()
		metrics.Count("request.events.size", len(message), "")
		payload := &de.EventRequest{}
		err = proto.Unmarshal(message, payload)
		if err != nil {
			logger.Error(fmt.Sprintf("[websocket.Handler] Reading message failed. %v  User ID: %s ", err, UserID))
			metrics.Count("batches.read", 1, "status=failed,reason=serde")
			badrequest := createBadrequestResponse(err)
			conn.WriteMessage(websocket.BinaryMessage, badrequest)
			continue
		}
		metrics.Count("batches.read", 1, "status=success")
		metrics.Count("request.events.count", len(payload.Events), "")

		wsHandler.bufferChannel <- EventsBatch{
			UserID:       UserID,
			EventReq:     payload,
			TimeConsumed: timeConsumed,
			TimePushed:   (time.Now()),
		}

		resp := createSuccessResponse(payload.ReqGuid)
		success, _ := proto.Marshal(resp)
		conn.WriteMessage(websocket.BinaryMessage, success)
	}
}

func calculateSessionTime(userID string, connectedAt time.Time) {
	connectionTime := time.Now().Sub(connectedAt)
	logger.Debug(fmt.Sprintf("[websocket.calculateSessionTime] UserID: %s, total time connected in minutes: %v", userID, connectionTime.Minutes()))
	metrics.Timing("users.session.time", connectionTime.Milliseconds(), "")
}

func setUpControlHandlers(conn *websocket.Conn, UserID string,
	PongWaitInterval time.Duration, WriteWaitInterval time.Duration) {
	//expects the client to send a ping, mark this channel as idle timed out post the deadline
	conn.SetReadDeadline(time.Now().Add(PongWaitInterval))
	conn.SetPongHandler(func(string) error {
		// extends the read deadline since we have received this pong on this channel
		conn.SetReadDeadline(time.Now().Add(PongWaitInterval))
		return nil
	})

	conn.SetPingHandler(func(s string) error {
		logger.Debug(fmt.Sprintf("Client connection with UserID: %s Pinged", UserID))
		if err := conn.WriteControl(websocket.PongMessage, []byte(s), time.Now().Add(WriteWaitInterval)); err != nil {
			metrics.Count("server.pong.failed", 1, "")
			logger.Debug(fmt.Sprintf("Failed to send pong event: %s UserID: %s", err, UserID))
		}
		return nil
	})
}
