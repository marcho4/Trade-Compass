package domain

import (
	"errors"
	"strconv"
)

type StockInfo struct {
	Ticker         string `json:"ticker"`
	NumberOfShares int    `json:"numberOfShares"`
	Name           string `json:"name"`
}

func ParseStockInfo(data [][]string) (*StockInfo, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	stockInfo := &StockInfo{}
	for _, row := range data {
		length := len(row)
		if length < 3 {
			return nil, errors.New("data is invalid")
		}
		switch row[0] {
		case "SECID":
			stockInfo.Ticker = row[2]
		case "NAME":
			stockInfo.Name = row[2]
		case "ISSUESIZE":
			numberOfShares, err := strconv.Atoi(row[2])
			if err != nil {
				return nil, err	
			}
			stockInfo.NumberOfShares = numberOfShares
		}
	}
	return stockInfo, nil
}
