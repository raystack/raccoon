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
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

//HandlerWSEvents handles the upgrade and the events sent by the peers
func (wsHandler *Handler) HandlerWSEvents(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("User-ID") // this should change with Proxy supplied header field
	connectedTime := time.Now()
	logger.Info(fmt.Sprintf("UserID %s connected at %v", userID, connectedTime))
	conn, err := wsHandler.websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("[websocket.Handler] Failed to upgrade connection: %v", err))
		return
	}
	defer conn.Close()
	defer calculateSessionTime(userID, connectedTime)

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
			resp := createUnknownrequestResponse(err)
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
