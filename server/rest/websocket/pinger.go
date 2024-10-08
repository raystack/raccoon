package websocket

import (
	"fmt"
	"time"

	"github.com/raystack/raccoon/core/identification"
	"github.com/raystack/raccoon/pkg/logger"
	"github.com/raystack/raccoon/pkg/metrics"
	"github.com/raystack/raccoon/server/rest/websocket/connection"
)

// Pinger is worker that pings the connected peers based on ping interval.
func Pinger(c chan connection.Conn, size int, PingInterval time.Duration, WriteWaitInterval time.Duration) {
	for range size {
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
							metrics.Increment("server_ping_failure_total", map[string]string{"conn_group": identifier.Group})
							delete(cSet, identifier)
						}
					}
				}
			}
		}()
	}
}
