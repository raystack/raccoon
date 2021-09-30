package connection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionPerGroup(t *testing.T) {
	t.Run("Should return all the group on the table with the count", func(t *testing.T) {
		table := NewTable(10)
		table.Store(Identifer{ID: "user1", Group: "group1"})
		table.Store(Identifer{ID: "user2", Group: "group1"})
		table.Store(Identifer{ID: "user3", Group: "group1"})
		table.Store(Identifer{ID: "user1", Group: "group2"})
		table.Store(Identifer{ID: "user2", Group: "group2"})
		assert.Equal(t, map[string]int{"group1": 3, "group2": 2}, table.TotalConnectionPerGroup())
	})
}

func TestStore(t *testing.T) {
	t.Run("Should store new connection", func(t *testing.T) {
		table := NewTable(10)
		table.Store(Identifer{ID: "user1", Group: ""})
		assert.True(t, table.Exists(Identifer{ID: "user1"}))
	})

	t.Run("Should remove connection when identifier match", func(t *testing.T) {
		table := NewTable(10)
		table.Store(Identifer{ID: "user1", Group: ""})
		table.Remove(Identifer{ID: "user1", Group: ""})
		assert.False(t, table.Exists(Identifer{ID: "user1", Group: ""}))
	})
}
