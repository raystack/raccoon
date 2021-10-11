package connection

import (
	"errors"
	"sync"
)

type Table struct {
	m       *sync.RWMutex
	connMap map[Identifier]struct{}
	counter map[string]int
	maxUser int
}

func NewTable(maxUser int) *Table {
	return &Table{
		m:       &sync.RWMutex{},
		connMap: make(map[Identifier]struct{}),
		maxUser: maxUser,
		counter: make(map[string]int),
	}
}

func (t *Table) Exists(c Identifier) bool {
	t.m.Lock()
	defer t.m.Unlock()
	_, ok := t.connMap[c]
	return ok
}

func (t *Table) Store(c Identifier) error {
	t.m.Lock()
	defer t.m.Unlock()
	if len(t.connMap) >= t.maxUser {
		return errMaxConnectionReached
	}
	if _, ok := t.connMap[c]; ok == true {
		return errConnDuplicated
	}
	t.connMap[c] = struct{}{}
	t.counter[c.Group] = t.counter[c.Group] + 1
	return nil
}

func (t *Table) Remove(c Identifier) {
	t.m.Lock()
	defer t.m.Unlock()
	delete(t.connMap, c)
	t.counter[c.Group] = t.counter[c.Group] - 1
}

func (t *Table) TotalConnection() int {
	t.m.Lock()
	defer t.m.Unlock()
	return len(t.connMap)
}

func (t *Table) TotalConnectionPerGroup() map[string]int {
	return t.counter
}

var errMaxConnectionReached = errors.New("max connection reached")

var errConnDuplicated = errors.New("duplicated connection")
