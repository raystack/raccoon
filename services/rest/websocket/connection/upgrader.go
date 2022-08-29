package connection

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/odpf/raccoon/identification"
	"github.com/odpf/raccoon/logger"
	"github.com/odpf/raccoon/metrics"
	pb "github.com/odpf/raccoon/proto"
	"google.golang.org/protobuf/proto"
)

type Upgrader struct {
	gorillaUg         websocket.Upgrader
	Table             *Table
	pongWaitInterval  time.Duration
	writeWaitInterval time.Duration
	connIDHeader      string
	connGroupHeader   string
	connGroupDefault  string
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
	ConnGroupDefault  string
}

func NewUpgrader(conf UpgraderConfig) *Upgrader {
	var checkOriginFunc func(r *http.Request) bool
	if !conf.CheckOrigin {
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
		connGroupDefault:  conf.ConnGroupDefault,
	}
}

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request) (Conn, error) {
	identifier := u.newIdentifier(r.Header)
	logger.Debug(fmt.Sprintf("%s connected at %v", identifier, time.Now()))

	conn, err := u.gorillaUg.Upgrade(w, r, nil)
	if err != nil {
		metrics.Increment("user_connection_failure_total", fmt.Sprintf("reason=ugfailure,conn_group=%s", identifier.Group))
		return Conn{}, fmt.Errorf("failed to upgrade %s: %v", identifier, err)
	}
	err = u.Table.Store(identifier)
	if errors.Is(err, errConnDuplicated) {
		errMsg := fmt.Sprintf("%s: %s,", err.Error(), identifier)
		duplicateConnResp := createEmptyErrorResponse(pb.Code_CODE_MAX_USER_LIMIT_REACHED, errMsg)

		conn.WriteMessage(websocket.BinaryMessage, duplicateConnResp)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "duplicate connection: "+identifier.ID))
		metrics.Increment("user_connection_failure_total", fmt.Sprintf("reason=exists,conn_group=%s", identifier.Group))
		conn.Close()
		return Conn{}, fmt.Errorf("disconnecting connection %s: already connected", identifier)
	}
	if errors.Is(err, errMaxConnectionReached) {
		logger.Errorf("[websocket.Handler] Disconnecting %v, max connection reached", identifier)
		maxConnResp := createEmptyErrorResponse(pb.Code_CODE_MAX_CONNECTION_LIMIT_REACHED, err.Error())
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

func (u *Upgrader) setUpControlHandlers(conn *websocket.Conn, identifier identification.Identifier) {
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

func (u *Upgrader) newIdentifier(h http.Header) identification.Identifier {
	// If connGroupHeader is empty string. By default, it will always return an empty string as Group. This means the group is fallback to default value.
	var group = h.Get(u.connGroupHeader)
	if group == "" {
		group = u.connGroupDefault
	}
	return identification.Identifier{
		ID:    h.Get(u.connIDHeader),
		Group: group,
	}
}

func createEmptyErrorResponse(errCode pb.Code, errMsg string) []byte {
	resp := pb.SendEventResponse{
		Status:   pb.Status_STATUS_ERROR,
		Code:     errCode,
		SentTime: time.Now().Unix(),
		Reason:   errMsg,
		Data:     nil,
	}
	duplicateConnResp, _ := proto.Marshal(&resp)
	return duplicateConnResp
}
