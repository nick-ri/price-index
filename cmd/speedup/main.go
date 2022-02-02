package main

import (
	"context"
	"github.com/NickRI/btc_index/cmd"
	"github.com/NickRI/btc_index/internal/aggregates"
	"github.com/NickRI/btc_index/internal/controller"
	"github.com/NickRI/btc_index/internal/indexes"
	"github.com/NickRI/btc_index/internal/models"
	"log"
	"os"
	"os/signal"
	"time"
)

const XTimesFaster = 10

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	btcUsdIdx := indexes.NewMBOptimized(models.BTCUSD, models.FiatDecimals, cmd.AmountOfSources*XTimesFaster, time.Minute)

	writers := cmd.MakeWriters(cmd.AmountOfSources, models.BTCUSD, 39250.12, 0.0003, time.Second/XTimesFaster, time.Second/XTimesFaster)

	aggr := aggregates.NewEfficient(btcUsdIdx, writers...)

	ctrl := controller.NewConsoleController(time.Second, btcUsdIdx, aggr)

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
