package collection

import (
	"context"
	"time"
)

type CollectFunction func(ctx context.Context, req *CollectRequest) error

func (c CollectFunction) Collect(ctx context.Context, req *CollectRequest) error {
	return c(ctx, req)
}

func NewChannelCollector(c chan *EventsBatch) Collector {
	return CollectFunction(func(ctx context.Context, req *CollectRequest) error {
		e := &EventsBatch{
			ConnectionIdentifier: req.ConnectionIdentifier,
			EventRequest:         req.EventRequest,
			TimeConsumed:         req.TimeConsumed,
			TimePushed:           time.Now(),
		}
		c <- e
		return nil
	})
}
