package domain

type Candle struct {
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Value  float64 `json:"value"`
	Volume float64 `json:"volume"`
	Begin  string  `json:"begin"`
	End    string  `json:"end"`
}
