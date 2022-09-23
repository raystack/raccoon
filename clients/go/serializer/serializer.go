package serializer

import (
	"errors"

	"encoding/json"

	"google.golang.org/protobuf/proto"
)

// SerializerFunc defines a conversion for raccoon message to byte sequence.
type SerializerFunc func(interface{}) ([]byte, error)

var (
	// json raccoon message serializer
	JSON = json.Marshal

	// proto raccoon message serializer
	PROTO = func(m interface{}) ([]byte, error) {
		msg, ok := m.(proto.Message)
		if !ok {
			return nil, errors.New("unable to marshal non proto")
		}
		return proto.Marshal(msg)
	}
)
