package domain

import "time"

type Dividends struct {
	ID              int       `json:"id"`
	Ticker          string    `json:"ticker"`
	ExDividendDate  time.Time `json:"exDividendDate"`
	PaymentDate     time.Time `json:"paymentDate"`
	AmountPerShare  float64   `json:"amountPerShare"`
	DividendYield   *float64  `json:"dividendYield,omitempty"`
	PayoutRatio     *float64  `json:"payoutRatio,omitempty"`
	Currency        string    `json:"currency"`
}
