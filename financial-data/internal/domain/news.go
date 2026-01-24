package domain

import "time"

type News struct {
	ticker  string
	date    time.Time
	content string
	source  string
}
