package log

import (
	"encoding/json"

	"github.com/raystack/raccoon/logger"
	pb "github.com/raystack/raccoon/proto"
)

// Publisher publishes message to the standard logger
// This is intended for development use.
type Publisher struct{}

func (p Publisher) ProduceBulk(events []*pb.Event, connGroup string) error {
	for _, event := range events {
		if json.Valid(event.EventBytes) {
			logger.Infof(
				"\nLogPublisher:\n\tmessage_type: json\n\tevent_type: %s\n\tevent: %s",
				event.Type,
				event.EventBytes,
			)
			continue
		}
		logger.Info(event.EventBytes)
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
	return Publisher{}
}
