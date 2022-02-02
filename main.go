package main

import (
	"context"
	"github.com/NickRI/btc_index/internal"
	"github.com/NickRI/btc_index/internal/aggregates"
	"github.com/NickRI/btc_index/internal/controller"
	"github.com/NickRI/btc_index/internal/indexes"
	"github.com/NickRI/btc_index/internal/infrastructure/sources"
	"github.com/NickRI/btc_index/internal/models"
	"log"
	"os"
	"os/signal"
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

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	btcUsdIdx := indexes.NewMutexBased(models.BTCUSD, AmountOfSources, time.Minute)

	subs := MakeSubscribers(AmountOfSources, models.BTCUSD, 39250.12, 0.003, time.Second, time.Second)

	aggr := aggregates.NewChannelBased(btcUsdIdx)

	ctrl := controller.NewConsoleController(time.Second, btcUsdIdx, aggr, subs...)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			<-c
			cancel()
		}()
	}()

	if err := ctrl.Exec(ctx); err != nil {
		log.Fatalf("controller executions fails with: %v", err)
	}
}
