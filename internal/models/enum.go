package models

//go:generate enumer -type=Ticker -linecomment -sql -json
type Ticker int

const (
	BTCUSD Ticker = iota + 1 // BTC_USD
)
