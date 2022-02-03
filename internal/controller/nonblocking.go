package controller

import (
	"context"
	"github.com/NickRI/btc_index/internal"
	"log"
)

type nonblockingController struct {
	index      internal.Index
	aggregator internal.Aggregator
	subscibers []internal.PriceStreamSubscriber
}

var _ internal.Controller = (*nonblockingController)(nil)

func NewNonBlockingController(
	index internal.Index,
	aggregator internal.Aggregator,
) *consoleController {
	return &consoleController{
		index:      index,
		aggregator: aggregator,
	}
}

func (n *nonblockingController) Exec(ctx context.Context) error {
	go func() {
		if err := n.aggregator.ListenStream(ctx); err != nil {
			log.Printf("index.ListenStream: got error, %v", err)
		}
	}()

	return nil
}
