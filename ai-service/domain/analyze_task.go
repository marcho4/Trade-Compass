package domain

type AnalyzeTask struct {
	ID     string `json:"id"`
	Ticker string `json:"ticker"`
	Year   int    `json:"year"`
	Period string `json:"period"`
}
