package domain

import "time"

type RefreshToken struct {
	id         int64
	tokenHash  string
	userId     int64
	deviceInfo string
	createdAt  time.Time
	updatedAt  time.Time
	expiresAt  time.Time
}
