package cmd

import (
	"time"

	"github.com/NickRI/btc_index/internal"
	"github.com/NickRI/btc_index/internal/infrastructure/sources"
	"github.com/NickRI/btc_index/internal/models"
)

const AmountOfSources = 100

func MakeSubscribers(count int, ticker models.Ticker, basePrice, errChance float64, minDelay, maxDelay time.Duration) []internal.PriceStreamSubscriber {
	var sims = make([]internal.PriceStreamSubscriber, 0, count)

	for i := 0; i < count; i++ {
		sims = append(sims, sources.NewSim(ticker, basePrice, errChance, minDelay, maxDelay))
	}

	return sims
}

func MakeWriters(count int, ticker models.Ticker, basePrice, errChance float64, minDelay, maxDelay time.Duration) []internal.PirceWriter {
	var sims = make([]internal.PirceWriter, 0, count)

	for i := 0; i < count; i++ {
		sims = append(sims, sources.NewSim(ticker, basePrice, errChance, minDelay, maxDelay))
	}

	return sims
}
