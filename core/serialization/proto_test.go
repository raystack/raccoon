package serialization_test

import (
	"testing"

	"github.com/raystack/raccoon/core/serialization"
	pb "github.com/raystack/raccoon/proto"
	"github.com/stretchr/testify/assert"
)

func TestSerialiseProto(t *testing.T) {
	t.Run("should return an error if argument is a non-protobuf message", func(t *testing.T) {
		arg := struct{}{}
		_, err := serialization.SerializeProto(arg)
		assert.Equal(t, err, serialization.ErrInvalidProtoMessage)
	})
	t.Run("should serialize a proto message", func(t *testing.T) {
		v := &pb.SendEventRequest{}
		_, err := serialization.SerializeProto(v)
		assert.Nil(t, err)
	})
}
