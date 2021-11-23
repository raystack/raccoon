package collection

import (
	"context"
	"time"

	"raccoon/identification"
	pb "raccoon/proto"
)

type CollectRequest struct {
	ConnectionIdentifier *identification.Identifier
	TimeConsumed         time.Time
	TimePushed           time.Time
	*pb.EventRequest
}

type Collector interface {
	Collect(ctx context.Context, req *CollectRequest) error
}
