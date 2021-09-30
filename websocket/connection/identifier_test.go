package connection

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifier(t *testing.T) {
	t.Run("Should extract id and group from specified header", func(t *testing.T) {
		header := http.Header{}
		header.Set("x-user-id", "user1")
		header.Set("x-user-group", "viewer")
		i := NewConnIdentifier(header, "x-user-id", "x-user-group")
		assert.Equal(t, Identifer{ID: "user1", Group: "viewer"}, i)
	})
}
