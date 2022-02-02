package controller

import (
	"context"
	"fmt"
	"github.com/NickRI/btc_index/internal"
	"log"
	"time"
)

type consoleController struct {
	rng        time.Duration
	index      internal.Index
	aggregator internal.Aggregator
	subscibers []internal.PriceStreamSubscriber
}

var _ internal.Controller = (*consoleController)(nil)

func NewConsoleController(
	rng time.Duration,
	index internal.Index,
	aggregator internal.Aggregator,
	subscibers ...internal.PriceStreamSubscriber,
) *consoleController {
	return &consoleController{
		rng:        rng,
		index:      index,
		aggregator: aggregator,
		subscibers: subscibers,
	}
}

func (c *consoleController) Exec(ctx context.Context) error {
	ticker := c.index.GetTicker()

	go func() {
		if err := c.aggregator.ListenStream(ctx, c.subscibers...); err != nil {
			log.Printf("aggregator.ListenStream: got error, %v", err)
		}
	}()

	fmt.Printf("index\ttimestamp\tprice\t\tfair\n") // TODO: move to presenter or smthg

	for result := range internal.RangePriceEmitter(ctx, c.index, c.rng) {
		fmt.Printf(
			"%s\t%d\t%s\t%s\t\n",
			ticker,
			time.Now().Unix(),
			result.Price.StringFixed(2),
			result.Fairness.StringFixed(2),
		)
	}

	return nil
}
