package websocket

import (
	"time"

	"github.com/goto/raccoon/metrics"
	"github.com/goto/raccoon/serialization"
	"github.com/goto/raccoon/services/rest/websocket/connection"
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
		metrics.Timing("ack_event_rtt_ms", ackTim.Milliseconds(), "")

		tim := time.Since(c.TimeConsumed)
		if c.Err != nil {
			metrics.Timing("event_rtt_ms", tim.Milliseconds(), "")
			writeFailedResponse(c.Conn, c.serializer, c.MessageType, c.RequestGuid, c.Err)
			continue
		}

		metrics.Timing("event_rtt_ms", tim.Milliseconds(), "")
		writeSuccessResponse(c.Conn, c.serializer, c.MessageType, c.RequestGuid)
	}
}
