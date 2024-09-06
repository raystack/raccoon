package log

import (
	"fmt"
	"testing"

	raccoonv1 "github.com/raystack/raccoon/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

type emitterProbe struct {
	Messages []string
}

func (pe *emitterProbe) Emit(value string) {
	pe.Messages = append(pe.Messages, value)
}

func TestLogPublisher(t *testing.T) {

	t.Run("should return an error if payload type is not of protobuf or json", func(t *testing.T) {
		payload := []*raccoonv1.Event{
			{
				Type:       "unknown",
				EventBytes: []byte("]{}"),
			},
		}
		p := Publisher{
			emit: func(string) {}, // noop
		}
		err := p.ProduceBulk(payload, "")
		assert.Error(t, err)
	})
	t.Run("should emit json events correctly", func(t *testing.T) {
		payload := []*raccoonv1.Event{
			{
				Type:       "unknown",
				EventBytes: []byte(`{"key":"value"}`),
			},
		}
		em := &emitterProbe{}
		p := Publisher{
			emit: em.Emit,
		}
		err := p.ProduceBulk(payload, "")
		assert.NoError(t, err)
		assert.Len(t, em.Messages, 1)

		expected := fmt.Sprintf(
			"[LogPublisher] kind = %s, event_type = %s, event = %s",
			"json",
			"unknown",
			`{"key":"value"}`,
		)
		assert.Equal(t, expected, em.Messages[0])
	})
	t.Run("should emit protobuf events correctly", func(t *testing.T) {
		msg := &raccoonv1.TestEvent{
			Description: "test event",
			Count:       420,
			Tags:        []string{"log", "protobuf"},
		}
		bytes, err := proto.Marshal(msg)
		require.NoError(t, err)
		payload := []*raccoonv1.Event{
			{
				Type:       "unknown",
				EventBytes: bytes,
			},
		}
		em := &emitterProbe{}
		p := Publisher{
			emit: em.Emit,
		}
		err = p.ProduceBulk(payload, "")
		assert.NoError(t, err)
		assert.Len(t, em.Messages, 1)

		expected := fmt.Sprintf(
			"[LogPublisher] kind = %s, event_type = %s, event = %s",
			"protobuf",
			"unknown",
			`1:"test event" 2:420 3:"log" 3:"protobuf"`,
		)
		assert.Equal(t, expected, em.Messages[0])
	})
	t.Run("publisher should have the name 'log'", func(t *testing.T) {
		assert.Equal(t, "log", Publisher{}.Name())
	})
	t.Run("publisher should Close() without error", func(t *testing.T) {
		p := Publisher{
			emit: (&emitterProbe{}).Emit,
		}
		err := p.ProduceBulk([]*raccoonv1.Event{
			{
				EventBytes: []byte("{}"),
				Type:       "unknown",
			},
		}, "")
		assert.NoError(t, err)
		assert.NoError(t, p.Close())
	})

	t.Run("publisher should initialise without panic", func(t *testing.T) {
		_ = New()
	})
}
