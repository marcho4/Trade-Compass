package domain

type Company struct {
	ID       int    `json:"id,omitempty"`
	Ticker   string `json:"ticker"`
	SectorID int    `json:"sectorId"`
	LotSize  int    `json:"lotSize,omitempty"`
	CEO      string `json:"ceo,omitempty"`
}
