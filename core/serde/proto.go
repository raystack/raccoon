package serde

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

var (
	ErrInvalidProtoMessage = errors.New("invalld proto message")
)

func SerializeProto(m interface{}) ([]byte, error) {
	msg, ok := m.(proto.Message)
	if !ok {
		return nil, ErrInvalidProtoMessage
	}
	return proto.Marshal(msg)
}

func DeserializeProto(b []byte, i interface{}) error {
	msg, ok := i.(proto.Message)
	if !ok {
		return ErrInvalidProtoMessage
	}
	return proto.Unmarshal(b, msg)
}
