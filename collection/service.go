package collection

import (
	"context"
	"time"
)

type CollectFunction func(ctx context.Context, req *CollectRequest) error

func (c CollectFunction) Collect(ctx context.Context, req *CollectRequest) error {
	return c(ctx, req)
}

func NewChannelCollector(c chan *CollectRequest) Collector {
	return CollectFunction(func(ctx context.Context, req *CollectRequest) error {
		req.TimePushed = time.Now()
		c <- req
		return nil
	})
}
