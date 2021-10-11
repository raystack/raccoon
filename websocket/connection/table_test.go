package connection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionPerGroup(t *testing.T) {
	t.Run("Should return all the group on the table with the count", func(t *testing.T) {
		table := NewTable(10)
		table.Store(Identifier{ID: "user1", Group: "group1"})
		table.Store(Identifier{ID: "user2", Group: "group1"})
		table.Store(Identifier{ID: "user3", Group: "group1"})
		table.Store(Identifier{ID: "user1", Group: "group2"})
		table.Store(Identifier{ID: "user2", Group: "group2"})
		assert.Equal(t, map[string]int{"group1": 3, "group2": 2}, table.TotalConnectionPerGroup())
	})
}

func TestStore(t *testing.T) {
	t.Run("Should store new connection", func(t *testing.T) {
		table := NewTable(10)
		err := table.Store(Identifier{ID: "user1", Group: ""})
		assert.NoError(t, err)
		assert.True(t, table.Exists(Identifier{ID: "user1"}))
	})

	t.Run("Should return max connection reached error when connection is maxed", func(t *testing.T) {
		table := NewTable(0)
		err := table.Store(Identifier{ID: "user1", Group: ""})
		assert.Error(t, err, errMaxConnectionReached)
	})

	t.Run("Should return duplicated error when connection already exists", func(t *testing.T) {
		table := NewTable(2)
		err := table.Store(Identifier{ID: "user1", Group: ""})
		assert.NoError(t, err)
		err = table.Store(Identifier{ID: "user1", Group: ""})
		assert.Error(t, err, errConnDuplicated)
	})

	t.Run("Should remove connection when identifier match", func(t *testing.T) {
		table := NewTable(10)
		table.Store(Identifier{ID: "user1", Group: ""})
		table.Remove(Identifier{ID: "user1", Group: ""})
		assert.False(t, table.Exists(Identifier{ID: "user1", Group: ""}))
	})
}
