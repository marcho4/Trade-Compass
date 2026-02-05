package domain

type Company struct {
	ID       int    `json:"id,omitempty"`
	Ticker   string `json:"ticker"`
	Name     string `json:"name,omitempty"`
	SectorID int    `json:"sectorId"`
	LotSize  int    `json:"lotSize,omitempty"`
	CEO      string `json:"ceo,omitempty"`
}
