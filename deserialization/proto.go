package deserialization

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

var ErrInvalidProtoMessage = errors.New("invalld proto message")

type ProtoDeserilizer struct{}

func (d *ProtoDeserilizer) Deserialize(b []byte, i interface{}) error {
	msg, ok := i.(proto.Message)
	if !ok {
		return ErrInvalidProtoMessage
	}
	return proto.Unmarshal(b, msg)

}
