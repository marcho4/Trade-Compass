package entity

type NewsSeverity string

const (
	SeverityHigh   NewsSeverity = "high"
	SeverityMedium NewsSeverity = "medium"
	SeverityLow    NewsSeverity = "low"
)

type NewsImpact string

const (
	ImpactPositive NewsImpact = "positive"
	ImpactNegative NewsImpact = "negative"
	ImpactNeutral  NewsImpact = "neutral"
)

type NewsItem struct {
	News       string       `json:"news"`
	Date       string       `json:"date"`
	Source     string       `json:"source"`
	Severity   NewsSeverity `json:"severity"`
	ImpactType NewsImpact   `json:"impact_type"`
}

type NewsResponse struct {
	LatestNews    []NewsItem `json:"latest_news"`
	ImportantNews []NewsItem `json:"important_news"`
}
