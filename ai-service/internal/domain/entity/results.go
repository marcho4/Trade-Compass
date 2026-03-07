package entity

type ReportResults struct {
	Health    int `json:"health"`
	Growth    int `json:"growth"`
	Moat      int `json:"moat"`
	Dividends int `json:"dividends"`
	Value     int `json:"value"`
	Total     int `json:"total"`
}

type AvailablePeriod struct {
	Year   int `json:"year"`
	Period int `json:"period"`
}
