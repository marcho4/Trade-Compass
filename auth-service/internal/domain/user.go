package domain

import "time"

type User struct {
	id          int64
	name        string
	status      string
	lastLoginAt time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

func (u *User) GetID() int64 {
	return u.id
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) GetStatus() string {
	return u.status
}