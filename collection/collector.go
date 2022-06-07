package collection

import (
	"context"
	"time"

	"github.com/odpf/raccoon/identification"
	pb "github.com/odpf/raccoon/proto"
)

type Event struct {
	Type       string
	EventBytes []byte
}

type CollectRequest struct {
	ConnectionIdentifier identification.Identifier
	TimeConsumed         time.Time
	TimePushed           time.Time
	*pb.SendEventRequest
}

type Collector interface {
	Collect(ctx context.Context, req *CollectRequest) error
}
