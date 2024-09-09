package collector

import (
	"context"

	"github.com/raystack/raccoon/pkg/clock"
)

type ChannelCollector struct {
	ch    chan CollectRequest
	clock clock.Clock
}

func NewChannelCollector(c chan CollectRequest) Collector {
	return &ChannelCollector{
		ch:    c,
		clock: clock.Default,
	}
}

func (c *ChannelCollector) Collect(ctx context.Context, req *CollectRequest) error {
	req.TimePushed = c.clock.Now()
	c.ch <- *req
	return nil
}
