package connection

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/odpf/raccoon/identification"
	"github.com/odpf/raccoon/logger"
	"github.com/odpf/raccoon/metrics"
)

type Conn struct {
	Identifier  identification.Identifier
	conn        *websocket.Conn
	connectedAt time.Time
	closeHook   func(c Conn)
	*sync.Mutex
}

func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	return c.conn.ReadMessage()
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	return c.conn.WriteMessage(messageType, data)
}

func (c *Conn) Ping(writeWaitInterval time.Duration) error {
	return c.conn.WriteControl(websocket.PingMessage, []byte("--ping--"), time.Now().Add(writeWaitInterval))
}

func (c *Conn) Close() {
	c.conn.Close()
	c.calculateSessionTime()
	c.closeHook(*c)
}

func (c *Conn) calculateSessionTime() {
	connectionTime := time.Since(c.connectedAt)
	logger.Debugf("[websocket.calculateSessionTime] %s, total time connected in minutes: %v", c.Identifier, connectionTime.Minutes())
	metrics.Timing("user_session_duration_milliseconds", connectionTime.Milliseconds(), fmt.Sprintf("conn_group=%s", c.Identifier.Group))
}
