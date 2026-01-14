package domain

type StockInfo struct {
	ticker string `json:"ticker"`
	numberOfShares int `json:"numberOfShares"`
	beta float64 `json:"beta"`
	marketRiskPremium float64 `json:"marketRiskPremium"`
	riskFreeRate float64 `json:"riskFreeRate"`
}