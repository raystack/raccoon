package collection

import (
	"context"
	"time"

	"github.com/odpf/raccoon/identification"
)

type Event struct {
	Type       string
	EventBytes []byte
}

type CollectRequest struct {
	ConnectionIdentifier identification.Identifier
	TimeConsumed         time.Time
	TimePushed           time.Time
	SentTime             time.Time
	Events               []Event
}

func (c CollectRequest) GetEvents() []Event {
	return c.Events
}

type Collector interface {
	Collect(ctx context.Context, req *CollectRequest) error
}
