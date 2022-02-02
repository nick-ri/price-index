package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type TickerPrice struct {
	Ticker Ticker
	Time   time.Time
	Price  string
}

type Numbers []decimal.Decimal

func (tps Numbers) Avg() (avg decimal.Decimal) {
	switch len(tps) {
	case 0:
		return decimal.Zero
	case 1:
		return tps[0]
	}

	sum := decimal.Zero

	for _, tp := range tps {
		sum = sum.Add(tp)
	}

	return sum.Div(decimal.NewFromInt(int64(len(tps))))
}

type PriceResult struct {
	Prices   Numbers
	Fairness Numbers
}

const (
	FiatDecimals = 2
)
