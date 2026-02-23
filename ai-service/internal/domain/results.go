package domain

type ReportResults struct {
	Health    int `json:"health"`
	Growth    int `json:"growth"`
	Moat      int `json:"moat"`
	Dividends int `json:"dividends"`
	Value     int `json:"value"`
	Total     int `json:"total"`
}
