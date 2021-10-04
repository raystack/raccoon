package connection

import (
	"errors"
	"fmt"
	"net/http"
	"raccoon/logger"
	"raccoon/metrics"
	pb "raccoon/websocket/proto"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

type Upgrader struct {
	gorillaUg         websocket.Upgrader
	Table             *Table
	pongWaitInterval  time.Duration
	writeWaitInterval time.Duration
	connIDHeader      string
	connGroupHeader   string
}

type UpgraderConfig struct {
	ReadBufferSize    int
	WriteBufferSize   int
	CheckOrigin       bool
	MaxUser           int
	PongWaitInterval  time.Duration
	WriteWaitInterval time.Duration
	ConnIDHeader      string
	ConnGroupHeader   string
}

func NewUpgrader(conf UpgraderConfig) *Upgrader {
	var checkOriginFunc func(r *http.Request) bool
	if conf.CheckOrigin == false {
		checkOriginFunc = func(r *http.Request) bool {
			return true
		}
	}
	return &Upgrader{
		gorillaUg: websocket.Upgrader{
			ReadBufferSize:  conf.ReadBufferSize,
			WriteBufferSize: conf.WriteBufferSize,
			CheckOrigin:     checkOriginFunc,
		},
		Table:             NewTable(conf.MaxUser),
		pongWaitInterval:  conf.PongWaitInterval,
		writeWaitInterval: conf.WriteWaitInterval,
		connIDHeader:      conf.ConnIDHeader,
		connGroupHeader:   conf.ConnGroupHeader,
	}
}

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request) (Conn, error) {
	identifier := Identifer{
		ID: r.Header.Get(u.connIDHeader),
		// If connGroupHeader is empty string. By default, it will always have an empty string as Group. This means uniqueness only depends on ID.
		Group: r.Header.Get(u.connGroupHeader),
	}
	logger.Debug(fmt.Sprintf("%s connected at %v", identifier, time.Now()))

	conn, err := u.gorillaUg.Upgrade(w, r, nil)
	if err != nil {
		metrics.Increment("user_connection_failure_total", fmt.Sprintf("reason=ugfailure,conn_group=%s", identifier.Group))
		return Conn{}, fmt.Errorf("failed to upgrade %s: %v", identifier, err)
	}
	err = u.Table.Store(identifier)
	if errors.Is(err, errConnDuplicated) {
		duplicateConnResp := createEmptyErrorResponse(pb.Code_MAX_USER_LIMIT_REACHED)

		conn.WriteMessage(websocket.BinaryMessage, duplicateConnResp)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "Duplicate connection"))
		metrics.Increment("user_connection_failure_total", fmt.Sprintf("reason=exists,conn_group=%s", identifier.Group))
		conn.Close()
		return Conn{}, fmt.Errorf("disconnecting %s: already connected", identifier)
	}
	if errors.Is(err, errMaxConnectionReached) {
		logger.Errorf("[websocket.Handler] Disconnecting %v, max connection reached", identifier)
		maxConnResp := createEmptyErrorResponse(pb.Code_MAX_CONNECTION_LIMIT_REACHED)
		conn.WriteMessage(websocket.BinaryMessage, maxConnResp)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "Max connection reached"))
		metrics.Increment("user_connection_failure_total", fmt.Sprintf("reason=serverlimit,conn_group=%s", identifier.Group))
		conn.Close()
		return Conn{}, fmt.Errorf("max connection reached")
	}

	u.setUpControlHandlers(conn, identifier)
	metrics.Increment("user_connection_success_total", fmt.Sprintf("conn_group=%s", identifier.Group))

	return Conn{
		Identifier:  identifier,
		conn:        conn,
		connectedAt: time.Now(),
		closeHook: func(c Conn) {
			u.Table.Remove(c.Identifier)
		}}, nil
}

func (u *Upgrader) setUpControlHandlers(conn *websocket.Conn, identifier Identifer) {
	//expects the client to send a ping, mark this channel as idle timed out post the deadline
	conn.SetReadDeadline(time.Now().Add(u.pongWaitInterval))
	conn.SetPongHandler(func(string) error {
		// extends the read deadline since we have received this pong on this channel
		conn.SetReadDeadline(time.Now().Add(u.pongWaitInterval))
		return nil
	})

	conn.SetPingHandler(func(s string) error {
		logger.Debug(fmt.Sprintf("Client %s pinged", identifier))
		if err := conn.WriteControl(websocket.PongMessage, []byte(s), time.Now().Add(u.writeWaitInterval)); err != nil {
			metrics.Increment("server_pong_failure_total", fmt.Sprintf("conn_group=%s", identifier.Group))
			logger.Debug(fmt.Sprintf("Failed to send pong event %s: %v", identifier, err))
		}
		return nil
	})
}

func createEmptyErrorResponse(errCode pb.Code) []byte {
	resp := pb.EventResponse{
		Status:   pb.Status_ERROR,
		Code:     errCode,
		SentTime: time.Now().Unix(),
		Reason:   "",
		Data:     nil,
	}
	duplicateConnResp, _ := proto.Marshal(&resp)
	return duplicateConnResp
}
