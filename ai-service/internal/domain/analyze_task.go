package domain

type Task struct {
	Ticker    string   `json:"ticker"`
	Year      int      `json:"year"`
	Period    string   `json:"period"`
	ReportURL string   `json:"report_url"`
	Type      TaskType `json:"type"`
}

type TaskType string

const (
	Analyze TaskType = "analyze"
	Extract TaskType = "extract"
)
