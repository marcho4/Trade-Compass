package entity

type Task struct {
	Id             string   `json:"id"`
	Ticker         string   `json:"ticker"`
	Year           int      `json:"year,omitempty"`
	Period         string   `json:"period,omitempty"`
	ReportURL      string   `json:"report_url,omitempty"`
	Type           TaskType `json:"type"`
	PendingCount   int      `json:"pending_count,omitempty"`
	ShouldContinue *bool    `json:"should_continue,omitempty"`
}

type TaskType string

const (
	Analyze              TaskType = "analyze"
	Extract              TaskType = "extract"
	ExtractResult        TaskType = "extract-result"
	BusinessResearch     TaskType = "business-research"
	NewsResearch         TaskType = "news-research"
	RiskAndGrowth        TaskType = "risk-and-growth"
	RiskAndGrowthExpect  TaskType = "risk-and-growth-expect"
	RiskAndGrowthSuccess TaskType = "risk-and-growth-success"
	RawDataExpect        TaskType = "raw-data-expect"
	RawDataSuccess       TaskType = "raw-data-success"
	GenerateScenarios    TaskType = "generate-scenarios"
	CalculateDCF         TaskType = "calculate-dcf"
)
