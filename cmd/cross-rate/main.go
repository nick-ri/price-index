package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/NickRI/btc_index/cmd"
	"github.com/NickRI/btc_index/internal"
	"github.com/NickRI/btc_index/internal/aggregates"
	"github.com/NickRI/btc_index/internal/controller"
	"github.com/NickRI/btc_index/internal/indexes"
	"github.com/NickRI/btc_index/internal/indexes/common"
	"github.com/NickRI/btc_index/internal/models"
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

	mainWriters := cmd.MakeWriters(
		cmd.AmountOfSources/10, // small amount for base currency to semulate bad fairness
		models.BTCUSD,
		39250.12,
		0.0001,
		time.Second/XTimesFaster,
		time.Second/XTimesFaster*2,
	)

	agg := aggregates.NewEfficient(crossBtcUsdIdx, mainWriters...)
	ctrl := controller.NewConsoleController(time.Minute, crossBtcUsdIdx, agg)

	ethUsdWriters := cmd.MakeWriters(
		cmd.AmountOfSources,
		models.ETHUSD,
		2669.80,
		0.0001,
		time.Second/XTimesFaster,
		time.Second/XTimesFaster*2,
	)

	ethBtcWriters := cmd.MakeWriters(
		cmd.AmountOfSources,
		models.ETHBTC,
		0.06,
		0.0001,
		time.Second/XTimesFaster,
		time.Second/XTimesFaster*2,
	)

	nonBlockingCtrls = append(nonBlockingCtrls,
		controller.NewNonBlockingController(ethUsdIds,
			aggregates.NewEfficient(ethUsdIds, ethUsdWriters...),
		),
		controller.NewNonBlockingController(ethBtcIds,
			aggregates.NewEfficient(ethBtcIds, ethBtcWriters...),
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
