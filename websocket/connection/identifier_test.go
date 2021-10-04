package connection

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifier(t *testing.T) {
	t.Run("Should extract id and group from specified header", func(t *testing.T) {
		header := http.Header{}
		header.Set("X-User-ID", "user1")
		header.Set("X-User-Group", "viewer")
		i := NewConnIdentifier(header, "X-User-ID", "X-User-Group")
		assert.Equal(t, Identifer{ID: "user1", Group: "viewer"}, i)
	})
}
