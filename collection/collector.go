package collection

import (
	"context"
	"time"

	pb "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/raccoon/v1beta1"
	"github.com/goto/raccoon/identification"
)

type AckFunc func(err error)

type CollectRequest struct {
	ConnectionIdentifier identification.Identifier
	TimeConsumed         time.Time
	TimePushed           time.Time
	AckFunc
	*pb.SendEventRequest
}

type Collector interface {
	Collect(ctx context.Context, req *CollectRequest) error
}
