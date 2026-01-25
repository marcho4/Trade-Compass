package domain

import "time"

type CBRate struct {
	Date time.Time `json:"date"`
	Rate float64   `json:"rate"`
}
