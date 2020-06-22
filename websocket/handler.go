package websocket

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"net/http"
	"raccoon/logger"
	"source.golabs.io/mobile/clickstream-go-proto/gojek/clickstream/de"
	"time"

	"github.com/gorilla/websocket"
)

type Handler struct {
	websocketUpgrader websocket.Upgrader
	bufferChannel     chan []*de.CSEventMessage
	user              *User
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

//HandlerWSEvents handles the upgrade and the events sent by the peers
func (wsHandler *Handler) HandlerWSEvents(w http.ResponseWriter, r *http.Request) {
	GoID := r.Header.Get("GO-User-ID")
	connectedTime := time.Now()
	logger.Info(fmt.Sprintf("GO-User-ID %s connected at %v", GoID, connectedTime))
	conn, err := wsHandler.websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("[websocket.Handler] Failed to upgrade connection: %v", err))
		return
	}
	defer conn.Close()
	defer calculateSessionTime(GoID, connectedTime)
	if wsHandler.user.Exists(GoID) {
		logger.Errorf("[websocket.Handler] Disconnecting %v, already connected", GoID)
		duplicateConnResp := createEmptyErrorResponse(de.Code_MAX_USER_LIMIT_REACHED)
		conn.WriteMessage(websocket.BinaryMessage, duplicateConnResp)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "Duplicated connection"))
		return
	}
	if wsHandler.user.HasReachedLimit() {
		logger.Errorf("[websocket.Handler] Disconnecting %v, max connection reached", GoID)
		maxConnResp := createEmptyErrorResponse(de.Code_MAX_CONNECTION_LIMIT_REACHED)
		conn.WriteMessage(websocket.BinaryMessage, maxConnResp)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "Max connection reached"))
		return
	}
	wsHandler.user.Store(GoID)
	defer wsHandler.user.Remove(GoID)

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
			logger.Error(fmt.Sprintf("[websocket.Handler] Reading message failed. Unknown failure: %v", err)) //no connection issue here
			break
		}
		request := &de.EventRequest{}
		err = proto.Unmarshal(message, request)
		if err != nil {
			logger.Error(fmt.Sprintf("[websocket.Handler] Reading message failed. %v", err))
			resp := createBadrequestResponse(err)
			unknownRequest, _ := proto.Marshal(&resp)
			conn.WriteMessage(websocket.BinaryMessage, unknownRequest)
			break
		}

		wsHandler.bufferChannel <- request.GetData()

		resp := createSuccessResponse(*request)
		success, _ := proto.Marshal(&resp)
		conn.WriteMessage(websocket.BinaryMessage, success)
	}
	/**
	* 1. @TODO - fetch user details from the header
	* 2. Verify if the user has connections already
	* 3. verify max connections for this server - How to respond to thr user in this case?
	* 4. Upgrade the connection
	* 5. Add this user-id -> connection mapping
	* 6. Add ping/pong handlers on this connection, readtimeout deadline
	* 6. Handle the message and send it to the events-channel - For now, as a go routine, deserialize protos
	* 7. Remove connection/user at the end of this function
	 */
}

func calculateSessionTime(userID string, connectedAt time.Time) {
	connectionTime := time.Now().Sub(connectedAt)
	logger.Info(fmt.Sprintf("[websocket.Handler] UserID: %s, total time connected in minutes: %v", userID, connectionTime.Minutes()))
	//@TODO - send this as metrics
}
