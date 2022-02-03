package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/NickRI/btc_index/internal"
)

type consoleController struct {
	rng        time.Duration
	index      internal.Index
	aggregator internal.Aggregator
}

var _ internal.Controller = (*consoleController)(nil)

func NewConsoleController(
	rng time.Duration,
	index internal.Index,
	aggregator internal.Aggregator,
) *consoleController {
	return &consoleController{
		rng:        rng,
		index:      index,
		aggregator: aggregator,
	}
}

func (c *consoleController) Exec(ctx context.Context) error {
	ticker := c.index.GetTicker()

	go func() {
		if err := c.aggregator.ListenStream(ctx); err != nil {
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
