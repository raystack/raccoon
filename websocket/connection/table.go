package connection

import (
	"sync"
)

type Table struct {
	m       sync.Mutex
	connMap map[Identifer]Identifer
	maxUser int
}

func NewTable(maxUser int) *Table {
	return &Table{
		m:       sync.Mutex{},
		connMap: make(map[Identifer]Identifer),
		maxUser: maxUser,
	}
}

func (u *Table) Exists(c Identifer) bool {
	u.m.Lock()
	defer u.m.Unlock()
	_, ok := u.connMap[c]
	return ok
}

func (u *Table) Store(c Identifer) {
	u.m.Lock()
	defer u.m.Unlock()
	u.connMap[c] = c
}

func (u *Table) Remove(c Identifer) {
	u.m.Lock()
	defer u.m.Unlock()
	delete(u.connMap, c)
}

func (u *Table) HasReachedLimit() bool {
	return u.TotalConnection() >= u.maxUser
}

func (u *Table) TotalConnection() int {
	u.m.Lock()
	defer u.m.Unlock()
	return len(u.connMap)
}
