package cmd

import (
	"github.com/NickRI/btc_index/internal"
	"github.com/NickRI/btc_index/internal/infrastructure/sources"
	"github.com/NickRI/btc_index/internal/models"
	"time"
)

const AmountOfSources = 100

func MakeSubscribers(count int, ticker models.Ticker, basePrice, errChance float64, minDelay, maxDelay time.Duration) []internal.PriceStreamSubscriber {
	var sims = make([]internal.PriceStreamSubscriber, 0, count)

	for i := 0; i < count; i++ {
		sims = append(sims, sources.NewSim(ticker, basePrice, errChance, minDelay, maxDelay))
	}

	return sims
}
