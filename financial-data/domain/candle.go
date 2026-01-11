package domain

type Candle struct {
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Value  float64
	Volume float64
	Begin  string
	End    string
}

func ParseCandles(candles [][]any) []Candle {
	result := make([]Candle, 0, len(candles))
	for _, x := range candles {
		if len(x) != 8 {
			continue
		}
		c := Candle{
			Open:   x[0].(float64),
			Close:  x[1].(float64),
			High:   x[2].(float64),
			Low:    x[3].(float64),
			Value:  x[4].(float64),
			Volume: x[5].(float64),
			Begin:  x[6].(string),
			End:    x[7].(string),
		}
		result = append(result, c)
	}
	return result
}
