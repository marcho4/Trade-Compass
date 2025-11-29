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
