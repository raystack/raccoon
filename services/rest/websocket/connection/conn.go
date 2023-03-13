package connection

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/goto/raccoon/identification"
	"github.com/goto/raccoon/logger"
	"github.com/goto/raccoon/metrics"
)

type Conn struct {
	Identifier  identification.Identifier
	conn        *websocket.Conn
	connectedAt time.Time
	closeHook   func(c Conn)
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
	if err := c.conn.Close(); err != nil {
		logger.Errorf("[Connection Error] %v", err)
		metrics.Increment("conn_close_err_count", "")
	}
	c.calculateSessionTime()
	c.closeHook(*c)
}

func (c *Conn) calculateSessionTime() {
	connectionTime := time.Since(c.connectedAt)
	logger.Debugf("[websocket.calculateSessionTime] %s, total time connected in minutes: %v", c.Identifier, connectionTime.Minutes())
	metrics.Timing("user_session_duration_milliseconds", connectionTime.Milliseconds(), fmt.Sprintf("conn_group=%s", c.Identifier.Group))
}
