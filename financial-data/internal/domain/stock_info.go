package domain

type StockInfo struct {
	ticker         string `json:"ticker"`
	numberOfShares int    `json:"numberOfShares"`
	stockMarket    string `json:"stockMarket"`
}
