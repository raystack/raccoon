package kinesis

import pb "github.com/raystack/raccoon/proto"

type Publisher struct{}

func (p *Publisher) ProduceBulk(events []*pb.Event, connGroup string) error {
	return nil
}

func New() *Publisher {
	return &Publisher{}
}
