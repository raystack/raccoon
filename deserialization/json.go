package deserialization

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func DeserializeJSON(b []byte, i interface{}) error {
	message, ok := i.(proto.Message)
	if !ok {
		return fmt.Errorf("expected a valid json which could be decoded into a protobuf")
	}
	return protojson.Unmarshal(b, message)
}
