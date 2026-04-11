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

type DependencyNewsItem struct {
	Dependency string       `json:"dependency"`
	News       string       `json:"news"`
	Date       string       `json:"date"`
	Source     string       `json:"source"`
	Severity   NewsSeverity `json:"severity"`
	ImpactType NewsImpact   `json:"impact_type"`
}

type NewsResponse struct {
	LatestNews               []NewsItem           `json:"latest_news"`
	HistoricalEvents         []NewsItem           `json:"historical_events"`
	UpcomingCompanyEvents    []NewsItem           `json:"upcoming_company_events"`
	UpcomingDependencyEvents []DependencyNewsItem `json:"upcoming_dependency_events"`
	PastDependencyEvents     []DependencyNewsItem `json:"past_dependency_events"`
}
