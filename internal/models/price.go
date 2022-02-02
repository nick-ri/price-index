package models

import (
	"time"
)

type TickerPrice struct {
	Ticker Ticker
	Time   time.Time
	Price  string
}
