package main

import (
	"context"
	"github.com/NickRI/btc_index/cmd"
	"github.com/NickRI/btc_index/internal"
	"github.com/NickRI/btc_index/internal/aggregates"
	"github.com/NickRI/btc_index/internal/controller"
	"github.com/NickRI/btc_index/internal/indexes"
	"github.com/NickRI/btc_index/internal/indexes/common"
	"github.com/NickRI/btc_index/internal/models"
	"log"
	"os"
	"os/signal"
	"time"
)

const XTimesFaster = 10

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	var (
		indexList        common.IndexList
		nonBlockingCtrls []internal.Controller
	)

	btcUsdIdx := indexes.NewLFOptimized(models.BTCUSD, models.FiatDecimals, cmd.AmountOfSources*XTimesFaster, time.Minute)
	ethUsdIds := indexes.NewLFOptimized(models.ETHUSD, models.FiatDecimals, cmd.AmountOfSources*XTimesFaster, time.Minute)
	ethBtcIds := indexes.NewLFOptimized(models.ETHBTC, models.FiatDecimals, cmd.AmountOfSources*XTimesFaster, time.Minute)

	indexList.AddItems(btcUsdIdx, ethUsdIds, ethBtcIds)

	crossBtcUsdIdx := indexes.NewCrossRate(btcUsdIdx, 0.5, indexList)

	writers := cmd.MakeWriters(cmd.AmountOfSources, models.BTCUSD, 39250.12, 0.001, time.Second/XTimesFaster, time.Second/XTimesFaster*2)

	agg := aggregates.NewEfficient(crossBtcUsdIdx, writers...)

	ctrl := controller.NewConsoleController(time.Minute, crossBtcUsdIdx, agg)

	nonBlockingCtrls = append(nonBlockingCtrls,
		controller.NewNonBlockingController(ethUsdIds,
			aggregates.NewEfficient(ethUsdIds, writers...),
		),
		controller.NewNonBlockingController(ethBtcIds,
			aggregates.NewEfficient(ethBtcIds, writers...),
		),
	)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			<-c
			cancel()
		}()
	}()

	for _, nbCtrl := range nonBlockingCtrls {
		if err := nbCtrl.Exec(ctx); err != nil {
			log.Fatalf("non-blocking controller execution fails with: %v", err)
		}
	}

	if err := ctrl.Exec(ctx); err != nil {
		log.Fatalf("controller execution fails with: %v", err)
	}
}
