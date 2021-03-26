package websocket

import (
	"fmt"
	"raccoon/logger"
	"raccoon/metrics"
	"time"

	"github.com/gorilla/websocket"
)

type connection struct {
	uniqConnID string
	conn       *websocket.Conn
}

//Pinger is a worker groroutine that pings the connected peers based on ping interval.
func Pinger(c chan connection, size int, PingInterval time.Duration, WriteWaitInterval time.Duration) {
	for i := 0; i < size; i++ {
		go func() {
			cSet := make(map[string]*websocket.Conn)
			timer := time.NewTicker(PingInterval)
			for {
				select {
				case conn := <-c:
					cSet[conn.uniqConnID] = conn.conn
				case <-timer.C:
					for uniqConnID, conn := range cSet {
						logger.Debug(fmt.Sprintf("Pinging UniqConnID: %s ", uniqConnID))
						if err := conn.WriteControl(websocket.PingMessage, []byte("--ping--"), time.Now().Add(WriteWaitInterval)); err != nil {
							logger.Error(fmt.Sprintf("[websocket.pingPeer] - Failed to ping User: %s Error: %v", uniqConnID, err))
							metrics.Increment("server_ping_failure_total", "")
							delete(cSet, uniqConnID)
						}
					}
				}
			}
		}()
	}
}
