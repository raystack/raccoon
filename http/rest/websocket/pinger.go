package websocket

import (
	"fmt"
	"raccoon/http/rest/websocket/connection"
	"raccoon/identification"
	"raccoon/logger"
	"raccoon/metrics"
	"time"
)

//Pinger is worker that pings the connected peers based on ping interval.
func Pinger(c chan connection.Conn, size int, PingInterval time.Duration, WriteWaitInterval time.Duration) {
	for i := 0; i < size; i++ {
		go func() {
			cSet := make(map[identification.Identifier]connection.Conn)
			ticker := time.NewTicker(PingInterval)
			for {
				select {
				case conn := <-c:
					cSet[conn.Identifier] = conn
				case <-ticker.C:
					for identifier, conn := range cSet {
						logger.Debug(fmt.Sprintf("Pinging %s ", identifier))
						if err := conn.Ping(WriteWaitInterval); err != nil {
							logger.Error(fmt.Sprintf("[websocket.pingPeer] - Failed to ping %s: %v", identifier, err))
							metrics.Increment("server_ping_failure_total", fmt.Sprintf("conn_group=%s", identifier.Group))
							delete(cSet, identifier)
						}
					}
				}
			}
		}()
	}
}
