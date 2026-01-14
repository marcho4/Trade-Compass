package domain

import "time"

type Dividends struct {
	ticker        string
	date          time.Time
	lastBuyDate   time.Time
	dividendYield float64
	payoutRatio   float64
	currency      string
	amount        float64
}
