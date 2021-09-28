package connection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionPerType(t *testing.T) {
	t.Run("Should return all the type on the table with the count", func(t *testing.T) {
		table := NewTable(10)
		table.Store(Identifer{ID: "user1", Type: "type1"})
		table.Store(Identifer{ID: "user2", Type: "type1"})
		table.Store(Identifer{ID: "user3", Type: "type1"})
		table.Store(Identifer{ID: "user1", Type: "type2"})
		table.Store(Identifer{ID: "user2", Type: "type2"})
		assert.Equal(t, map[string]int{"type1": 3, "type2": 2}, table.ConnectionPerType())
	})
}

func TestStore(t *testing.T) {
	t.Run("Should store new connection", func(t *testing.T) {
		table := NewTable(10)
		table.Store(Identifer{ID: "user1", Type: ""})
		assert.True(t, table.Exists(Identifer{ID: "user1"}))
	})

	t.Run("Should remove connection when identifier match", func(t *testing.T) {
		table := NewTable(10)
		table.Store(Identifer{ID: "user1", Type: ""})
		table.Remove(Identifer{ID: "user1", Type: ""})
		assert.False(t, table.Exists(Identifer{ID: "user1", Type: ""}))
	})
}
