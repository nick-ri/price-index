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

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	btcUsdIdx := indexes.NewMutexBased(models.BTCUSD, cmd.AmountOfSources, time.Minute)

	subs := cmd.MakeSubscribers(cmd.AmountOfSources, models.BTCUSD, 39250.12, 0.003, time.Second, time.Second*2)

	aggr := aggregates.NewChannelBased(btcUsdIdx, subs...)

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
