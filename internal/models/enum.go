package models

import "strings"

//go:generate enumer -type=Ticker -linecomment -sql -json
type Ticker int

const (
	BTCUSD Ticker = iota + 1 // BTC_USD
	ETHUSD                   // ETH_USD
	ETHBTC                   // ETH_BTC
)

func (i Ticker) Base() string {
	return strings.Split(i.String(), "_")[0]
}

func (i Ticker) Secondary() string {
	return strings.Split(i.String(), "_")[1]
}
