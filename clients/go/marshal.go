package raccoon

import (
	"errors"

	"encoding/json"

	"google.golang.org/protobuf/proto"
)

type Message interface{}

var (
	// json clickstream marshaler
	JSON = json.Marshal

	// proto clickstream marshaler
	PROTO = func(m interface{}) ([]byte, error) {
		msg, ok := m.(proto.Message)
		if !ok {
			return nil, errors.New("unable to marshal non proto")
		}
		return proto.Marshal(msg)
	}
)

// Marshaler defines a conversion for clickstream message to byte sequence.
type Marshaler interface {
	Marshal(any Message) ([]byte, error)
}

type MarshalFunc func(interface{}) ([]byte, error)

// Marshal marshals "any" clickstream message into byte sequence.
func (m MarshalFunc) Marshal(any Message) ([]byte, error) {
	return m(any)
}
