package log

import (
	"cmp"
	"encoding/json"
	"fmt"

	"github.com/raystack/raccoon/pkg/logger"
	pb "github.com/raystack/raccoon/proto"
	"github.com/raystack/raccoon/publisher"
	"github.com/turtleDev/protoraw"
)

type logEmitter func(string)

// Publisher publishes message to the standard logger
// This is intended for development use.
type Publisher struct {
	emit logEmitter
}

func (p Publisher) ProduceBulk(events []*pb.Event, connGroup string) error {
	var errs []error
	for _, e := range events {
		var (
			typ   = e.Type
			kind  string
			event string
		)
		if json.Valid(e.EventBytes) {
			kind = "json"
			event = string(e.EventBytes)
		} else {
			kind = "protobuf"
			var err error
			event, err = protoraw.Decode(e.EventBytes)
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}
		line := fmt.Sprintf(
			"[LogPublisher] kind = %s, event_type = %s, event = %s",
			kind,
			typ,
			event,
		)
		p.emit(line)
	}
	if cmp.Or(errs...) != nil {
		return &publisher.BulkError{
			Errors: errs,
		}
	}
	return nil
}

func (p Publisher) Name() string {
	return "log"
}

func (p Publisher) Close() error {
	return nil
}

func New() Publisher {
	return Publisher{
		emit: func(line string) {
			logger.Info(line)
		},
	}
}
