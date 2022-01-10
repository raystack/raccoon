package serialization

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

var (
	ErrInvalidProtoMessage = errors.New("invalld proto message")
)

type ProtoSerilizer struct{}

func (p *ProtoSerilizer) Serialize(m interface{}) ([]byte, error) {
	msg, ok := m.(proto.Message)
	if !ok {
		return nil, ErrInvalidProtoMessage
	}
	return proto.Marshal(msg)
}
