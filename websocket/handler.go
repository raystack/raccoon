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
	bufferChannel     chan de.EventRequest
	user              *User
	PingInterval      time.Duration
	PongWaitInterval  time.Duration
	WriteWaitInterval time.Duration
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

//HandlerWSEvents handles the upgrade and the events sent by the peers
func (wsHandler *Handler) HandlerWSEvents(w http.ResponseWriter, r *http.Request) {
	GOUserID := r.Header.Get("GO-User-ID")
	connectedTime := time.Now()
	logger.Info(fmt.Sprintf("GO-User-ID %s connected at %v", GOUserID, connectedTime))
	conn, err := wsHandler.websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("[websocket.Handler] Failed to upgrade connection  User ID: %s : %v", GOUserID, err))
		metrics.Increment("users.disconnected", "reason=ugfailure")
		return
	}
	defer conn.Close()

	if wsHandler.user.Exists(GOUserID) {
		logger.Errorf("[websocket.Handler] Disconnecting %v, already connected", GOUserID)
		duplicateConnResp := createEmptyErrorResponse(de.Code_MAX_USER_LIMIT_REACHED)

		conn.WriteMessage(websocket.BinaryMessage, duplicateConnResp)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "Duplicate connection"))
		metrics.Count("user.disconnected", 1, "reason=exists")
		return
	}
	if wsHandler.user.HasReachedLimit() {
		logger.Errorf("[websocket.Handler] Disconnecting %v, max connection reached", GOUserID)
		maxConnResp := createEmptyErrorResponse(de.Code_MAX_CONNECTION_LIMIT_REACHED)
		conn.WriteMessage(websocket.BinaryMessage, maxConnResp)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "Max connection reached"))
		metrics.Count("user.disconnected", 1, "reason=serverlimit")
		return
	}
	wsHandler.user.Store(GOUserID)
	defer wsHandler.user.Remove(GOUserID)
	defer calculateSessionTime(GOUserID, connectedTime)

	setUpControlHandlers(conn, GOUserID, wsHandler.PingInterval, wsHandler.PongWaitInterval, wsHandler.WriteWaitInterval)
	metrics.Increment("user.connected", "")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
				websocket.CloseNoStatusReceived,
				websocket.CloseAbnormalClosure) {
				logger.Error(fmt.Sprintf("[websocket.Handler] Connection Closed Abruptly: %v", err))
				break
			}
			logger.Error(fmt.Sprintf("[websocket.Handler] Reading message failed. Unknown failure: %v  User ID: %s ", err, GOUserID)) //no connection issue here
			break
		}
		metrics.Count("request.events.bytes", len(message), "")
		payload := &de.EventRequest{}
		err = proto.Unmarshal(message, payload)
		if err != nil {
			logger.Error(fmt.Sprintf("[websocket.Handler] Reading message failed. %v  User ID: %s ", err, GOUserID))
			metrics.Increment("request.count", "type=bad")
			badrequest := createBadrequestResponse(err)
			conn.WriteMessage(websocket.BinaryMessage, badrequest)
			continue
		}
		metrics.Increment("request.count", "type=ok")

		wsHandler.bufferChannel <- *payload

		resp := createSuccessResponse(*payload)
		success, _ := proto.Marshal(&resp)
		conn.WriteMessage(websocket.BinaryMessage, success)
	}
}

func calculateSessionTime(userID string, connectedAt time.Time) {
	connectionTime := time.Now().Sub(connectedAt)
	logger.Info(fmt.Sprintf("[websocket.calculateSessionTime] UserID: %s, total time connected in minutes: %v", userID, connectionTime.Minutes()))
	metrics.Timing("users.session.time", connectionTime.Milliseconds(), "")
}

func setUpControlHandlers(conn *websocket.Conn, GOUserID string, PingInterval time.Duration,
	PongWaitInterval time.Duration, WriteWaitInterval time.Duration) {
	//expects the client to send a ping, mark this channel as idle timed out post the deadline
	conn.SetReadDeadline(time.Now().Add(PongWaitInterval))
	conn.SetPongHandler(func(string) error {
		// extends the read deadline since we have received this pong on this channel
		conn.SetReadDeadline(time.Now().Add(PongWaitInterval))
		return nil
	})

	conn.SetPingHandler(func(s string) error {
		logger.Debug(fmt.Sprintf("Client connection with User ID: %s Pinged", GOUserID))
		if err := conn.WriteControl(websocket.PongMessage, []byte(s), time.Now().Add(WriteWaitInterval)); err != nil {
			logger.Debug(fmt.Sprintf("Failed to send ping event: %s User: %s", err, GOUserID))
		}
		return nil
	})
	go pingPeer(GOUserID, conn, PingInterval, WriteWaitInterval)
}

func pingPeer(userID string, conn *websocket.Conn, PingInterval time.Duration, WriteWaitInterval time.Duration) {
	timer := time.NewTicker(PingInterval)
	defer func() {
		timer.Stop()
	}()

	for {
		<-timer.C
		logger.Debug(fmt.Sprintf("Pinging UserId: %s ", userID))
		if err := conn.WriteControl(websocket.PingMessage, []byte("--ping--"), time.Now().Add(WriteWaitInterval)); err != nil {
			logger.Error(fmt.Sprintf("[websocket.pingPeer] - Failed to ping User: %s Error: %v", userID, err))
			return
		}
	}
}
