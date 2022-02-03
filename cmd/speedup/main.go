package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/NickRI/price-index/cmd"
	"github.com/NickRI/price-index/internal/aggregates"
	"github.com/NickRI/price-index/internal/controller"
	"github.com/NickRI/price-index/internal/indexes"
	"github.com/NickRI/price-index/internal/models"
)

const XTimesFaster = 10

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	btcUsdIdx := indexes.NewMBOptimized(models.BTCUSD, models.FiatDecimals, cmd.AmountOfSources*XTimesFaster, time.Minute)

	writers := cmd.MakeWriters(cmd.AmountOfSources, models.BTCUSD, 39250.12, 0.0003, time.Second/XTimesFaster, time.Second/XTimesFaster*2)

	aggr := aggregates.NewEfficient(btcUsdIdx, writers...)

	ctrl := controller.NewConsoleController(time.Minute, btcUsdIdx, aggr)

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
