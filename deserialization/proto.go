package deserialization

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

var ErrInvalidProtoMessage = errors.New("invalld proto message")

func ProtoDeserilizer() Deserializer {
	return DeserializeFunc(func(b []byte, i interface{}) error {
		msg, ok := i.(proto.Message)
		if !ok {
			return ErrInvalidProtoMessage
		}
		return proto.Unmarshal(b, msg)
	})

}
