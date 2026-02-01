package domain

type CandlesApiResponse struct {
	Candles Candles `json:"candles"`
}

type Candles struct {
	Metadata map[string]ColumnMetadata `json:"metadata"`
	Columns  []string                  `json:"columns"`
	Data     [][]any                   `json:"data"`
}

type ColumnMetadata struct {
	Type    string `json:"type"`
	Bytes   int    `json:"bytes,omitempty"`
	MaxSize int    `json:"max_size,omitempty"`
}

type StockInfoApiResponse struct {
	Description Description `json:"description"`
}

type Description struct {
	Columns []string   `json:"columns"`
	Data    [][]string `json:"data"`
}
