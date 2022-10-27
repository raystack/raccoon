package connection

import (
	"errors"
	"sync"

	"github.com/odpf/raccoon/identification"
)

var (
	errMaxConnectionReached = errors.New("max connection reached")
	errConnDuplicated       = errors.New("duplicated connection")
)

type Table struct {
	m       *sync.RWMutex
	connMap map[identification.Identifier]map[string]struct{}
	counter map[string]int
	maxUser int
}

func NewTable(maxUser int) *Table {
	return &Table{
		m:       &sync.RWMutex{},
		connMap: make(map[identification.Identifier]map[string]struct{}),
		maxUser: maxUser,
		counter: make(map[string]int),
	}
}

func (t *Table) Exists(c identification.Identifier) bool {
	t.m.Lock()
	defer t.m.Unlock()
	_, ok := t.connMap[c]
	return ok
}

func (t *Table) Store(c identification.Identifier) error {
	t.m.Lock()
	defer t.m.Unlock()
	if len(t.connMap) >= t.maxUser {
		return errMaxConnectionReached
	}
	if _, ok := t.connMap[c]; ok {
		return errConnDuplicated
	}
	t.connMap[c] = make(map[string]struct{})
	t.counter[c.Group] = t.counter[c.Group] + 1
	return nil
}

func (t *Table) StoreBatch(c identification.Identifier, id string) {
	t.m.Lock()
	defer t.m.Unlock()
	if _, ok := t.connMap[c]; ok {
		t.connMap[c][id] = struct{}{}
	}
}

func (t *Table) HasBatch(c identification.Identifier, id string) bool {
	t.m.Lock()
	defer t.m.Unlock()
	_, ok := t.connMap[c][id]
	return ok
}

func (t *Table) Remove(c identification.Identifier) {
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
