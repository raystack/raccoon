package log

import (
	"github.com/raystack/raccoon/logger"
	pb "github.com/raystack/raccoon/proto"
)

// Publisher publishes message to the standard logger
// This is intended for development use.
type Publisher struct{}

func (p Publisher) ProduceBulk(events []*pb.Event, connGroup string) error {
	for _, event := range events {
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
