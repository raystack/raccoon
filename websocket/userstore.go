package websocket

import (
	"sync"
)

type User struct {
	m       sync.Mutex
	userMap map[ConnIdentifier]ConnIdentifier
	maxUser int
}

func NewUserStore(maxUser int) *User {
	return &User{
		m:       sync.Mutex{},
		userMap: make(map[ConnIdentifier]ConnIdentifier),
		maxUser: maxUser,
	}
}

func (u *User) Exists(c ConnIdentifier) bool {
	u.m.Lock()
	defer u.m.Unlock()
	_, ok := u.userMap[c]
	return ok
}

func (u *User) Store(c ConnIdentifier) {
	u.m.Lock()
	defer u.m.Unlock()
	u.userMap[c] = c
}

func (u *User) Remove(c ConnIdentifier) {
	u.m.Lock()
	defer u.m.Unlock()
	delete(u.userMap, c)
}

func (u *User) HasReachedLimit() bool {
	return u.TotalUsers() >= u.maxUser
}

func (u *User) TotalUsers() int {
	u.m.Lock()
	defer u.m.Unlock()
	return len(u.userMap)
}
