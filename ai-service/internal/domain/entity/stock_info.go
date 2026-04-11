package entity

type StockInfo struct {
	Ticker         string `json:"ticker"`
	NumberOfShares int    `json:"numberOfShares"`
	Name           string `json:"name"`
}
