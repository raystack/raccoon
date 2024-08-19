package log

import (
	"encoding/json"

	"github.com/raystack/raccoon/logger"
	pb "github.com/raystack/raccoon/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
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
		fdp := &descriptorpb.FileDescriptorProto{
			Name: proto.String("empty_message.proto"),
			MessageType: []*descriptorpb.DescriptorProto{
				&descriptorpb.DescriptorProto{
					Name: proto.String("EmptyMessage"),
				},
			},
		}
		fd, err := protodesc.NewFile(fdp, &protoregistry.Files{})
		if err != nil {
			// todo
			panic(err)
		}
		m := dynamicpb.NewMessage(fd.Messages().ByName("EmptyMessage"))
		proto.Unmarshal(event.EventBytes, m)
		logger.Info(m.String())
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
