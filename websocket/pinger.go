package websocket

import (
	"fmt"
	"raccoon/logger"
	"raccoon/metrics"
	"raccoon/websocket/connection"
	"time"
)

//Pinger is worker that pings the connected peers based on ping interval.
func Pinger(c chan connection.Conn, size int, PingInterval time.Duration, WriteWaitInterval time.Duration) {
	for i := 0; i < size; i++ {
		go func() {
			cSet := make(map[connection.Identifer]connection.Conn)
			timer := time.NewTicker(PingInterval)
			for {
				select {
				case conn := <-c:
					cSet[conn.Identifier] = conn
				case <-timer.C:
					for identifier, conn := range cSet {
						logger.Debug(fmt.Sprintf("Pinging %s ", identifier))
						if err := conn.Ping(WriteWaitInterval); err != nil {
							logger.Error(fmt.Sprintf("[websocket.pingPeer] - Failed to ping %s: %v", identifier, err))
							metrics.Increment("server_ping_failure_total", fmt.Sprintf("conn_type=%s", identifier.Type))
							delete(cSet, identifier)
						}
					}
				}
			}
		}()
	}
}
