package entity

type CompanyProfile struct {
	Ticker              string   `json:"ticker"`
	CompanyName         string   `json:"company_name"`
	Description         string   `json:"description"`
	ProductsAndServices []string `json:"products_and_services"`
	Markets             []Market `json:"markets"`
	KeyClients          string   `json:"key_clients"`
	BusinessModel       string   `json:"business_model"`
}

type Market struct {
	Market string `json:"market"`
	Role   string `json:"role"`
}

type RevenueSource struct {
	Ticker      string  `json:"ticker"`
	Segment     string  `json:"segment"`
	SharePct    float64 `json:"share_pct"`
	Approximate bool    `json:"approximate"`
	Description string  `json:"description"`
	Trend       string  `json:"trend"`
}

type CompanyDependency struct {
	Ticker      string `json:"ticker"`
	Factor      string `json:"factor"`
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

type BusinessResearchResponse struct {
	Ticker      string `json:"ticker"`
	CompanyName string `json:"company_name"`
	Profile     struct {
		Description         string   `json:"description"`
		ProductsAndServices []string `json:"products_and_services"`
		Markets             []Market `json:"markets"`
		KeyClients          string   `json:"key_clients"`
		BusinessModel       string   `json:"business_model"`
	} `json:"profile"`
	RevenueSources []struct {
		Segment     string  `json:"segment"`
		SharePct    float64 `json:"share_pct"`
		Approximate bool    `json:"approximate"`
		Description string  `json:"description"`
		Trend       string  `json:"trend"`
	} `json:"revenue_sources"`
	Dependencies []struct {
		Factor      string `json:"factor"`
		Type        string `json:"type"`
		Severity    string `json:"severity"`
		Description string `json:"description"`
	} `json:"dependencies"`
}

type BusinessResearchResult struct {
	Profile      CompanyProfile      `json:"profile"`
	Revenue      []RevenueSource     `json:"revenue_sources"`
	Dependencies []CompanyDependency `json:"dependencies"`
}
