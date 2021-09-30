package connection

import (
	"sync"
)

type Table struct {
	m       sync.Mutex
	connMap map[Identifer]Identifer
	counter map[string]int
	maxUser int
}

func NewTable(maxUser int) *Table {
	return &Table{
		m:       sync.Mutex{},
		connMap: make(map[Identifer]Identifer),
		maxUser: maxUser,
		counter: make(map[string]int),
	}
}

func (t *Table) Exists(c Identifer) bool {
	t.m.Lock()
	defer t.m.Unlock()
	_, ok := t.connMap[c]
	return ok
}

func (t *Table) Store(c Identifer) {
	t.m.Lock()
	defer t.m.Unlock()
	t.connMap[c] = c
	t.counter[c.Group] = t.counter[c.Group] + 1
}

func (t *Table) Remove(c Identifer) {
	t.m.Lock()
	defer t.m.Unlock()
	delete(t.connMap, c)
	t.counter[c.Group] = t.counter[c.Group] - 1
}

func (t *Table) HasReachedLimit() bool {
	return t.TotalConnection() >= t.maxUser
}

func (t *Table) TotalConnection() int {
	t.m.Lock()
	defer t.m.Unlock()
	return len(t.connMap)
}

func (t *Table) TotalConnectionPerGroup() map[string]int {
	return t.counter
}
