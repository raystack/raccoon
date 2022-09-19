package wire

import (
	"errors"

	"encoding/json"

	"google.golang.org/protobuf/proto"
)

// WireMarshaler defines a conversion between byte sequence and raccoon request payloads.
type WireMarshaler interface {

	// Marshal marshals "any" into byte sequence.
	Marshal(any interface{}) ([]byte, error)

	// Unmarshal unmarshals "data" into "any".
	// "any" must be a pointer value.
	Unmarshal(data []byte, any interface{}) error

	// ContentType returns the Content-Type which this marshaler is responsible for.
	ContentType() string
}

// JsonWire is a Marshaler which marshals/unmarshals into/from serialize json bytes
type JsonWire struct{}

// ContentType always returns "application/json".
func (*JsonWire) ContentType() string {
	return "application/json"
}

// Marshal marshals "any" into JSON
func (j *JsonWire) Marshal(any interface{}) ([]byte, error) {
	return json.Marshal(any)
}

// Unmarshal unmarshals JSON data into "any".
func (j *JsonWire) Unmarshal(data []byte, any interface{}) error {
	return json.Unmarshal(data, any)
}

// ProtoWire is a Marshaler which marshals/unmarshals into/from serialize proto bytes
type ProtoWire struct{}

// ContentType always returns "application/proto".
func (*ProtoWire) ContentType() string {
	return "application/proto"
}

// Marshal marshals "value" into Proto
func (*ProtoWire) Marshal(value interface{}) ([]byte, error) {
	message, ok := value.(proto.Message)
	if !ok {
		return nil, errors.New("unable to marshal non proto")
	}
	return proto.Marshal(message)
}

// Unmarshal unmarshals proto "data" into "value"
func (*ProtoWire) Unmarshal(data []byte, value interface{}) error {
	message, ok := value.(proto.Message)
	if !ok {
		return errors.New("unable to unmarshal non proto")
	}
	return proto.Unmarshal(data, message)
}
