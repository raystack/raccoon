package websocket

import (
	"time"

	"github.com/raystack/raccoon/core/serialization"
	"github.com/raystack/raccoon/pkg/metrics"
	"github.com/raystack/raccoon/server/rest/websocket/connection"
)

var AckChan = make(chan AckInfo)

type AckInfo struct {
	MessageType     int
	RequestGuid     string
	Err             error
	Conn            connection.Conn
	serializer      serialization.SerializeFunc
	TimeConsumed    time.Time
	AckTimeConsumed time.Time
}

func AckHandler(ch <-chan AckInfo) {
	for c := range ch {
		ackTim := time.Since(c.AckTimeConsumed)
		metrics.Histogram("ack_event_rtt_ms", ackTim.Milliseconds(), map[string]string{})

		tim := time.Since(c.TimeConsumed)
		if c.Err != nil {
			metrics.Histogram("event_rtt_ms", tim.Milliseconds(), map[string]string{})
			writeFailedResponse(c.Conn, c.serializer, c.MessageType, c.RequestGuid, c.Err)
			continue
		}

		metrics.Histogram("event_rtt_ms", tim.Milliseconds(), map[string]string{})
		writeSuccessResponse(c.Conn, c.serializer, c.MessageType, c.RequestGuid)
	}
}
