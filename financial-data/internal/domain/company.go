package domain

type Company struct {
	ID        int    `json:"id"`
	INN       string `json:"inn"`
	Ticker    string `json:"ticker"`
	Owner     string `json:"owner"`
	SectorID  int    `json:"sectorId"`
	LotSize   int    `json:"lotSize"`
	CEO       string `json:"ceo,omitempty"`
	Employees *int   `json:"employees,omitempty"`
}
