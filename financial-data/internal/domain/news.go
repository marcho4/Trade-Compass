package domain

import "time"

type News struct {
	ID       int       `json:"id"`
	Ticker   *string   `json:"ticker,omitempty"`
	SectorID *int      `json:"sectorId,omitempty"`
	Date     time.Time `json:"date"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Source   string    `json:"source"`
	URL      *string   `json:"url,omitempty"`
}
