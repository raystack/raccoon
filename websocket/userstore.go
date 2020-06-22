package websocket

import "sync"

type User struct {
	m       sync.Mutex
	userMap map[string]string
	maxUser int
}

func NewUserStore(maxUser int) *User {
	return &User{
		m:       sync.Mutex{},
		userMap: make(map[string]string),
		maxUser: maxUser,
	}
}

func (u *User) Exists(userId string) bool {
	u.m.Lock()
	defer u.m.Unlock()
	_, ok := u.userMap[userId]
	return ok
}

func (u *User) Store(userId string) {
	u.m.Lock()
	defer u.m.Unlock()
	u.userMap[userId] = userId
}

func (u *User) Remove(userId string) {
	u.m.Lock()
	defer u.m.Unlock()
	delete(u.userMap, userId)
}

func (u *User) HasReachedLimit() bool {
	return u.TotalUsers() >= u.maxUser
}

func (u *User) TotalUsers() int {
	u.m.Lock()
	defer u.m.Unlock()
	return len(u.userMap)
}
