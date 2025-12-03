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

func NewRefreshToken(id int64, tokenHash string, userID int64, deviceInfo string, createdAt, updatedAt, expiresAt time.Time) *RefreshToken {
	return &RefreshToken{
		id:         id,
		tokenHash:  tokenHash,
		userId:     userID,
		deviceInfo: deviceInfo,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
		expiresAt:  expiresAt,
	}
}
