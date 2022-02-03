package controller

import (
	"context"
	"log"

	"github.com/NickRI/price-index/internal"
)

type nonblockingController struct {
	index      internal.Index
	aggregator internal.Aggregator
}

var _ internal.Controller = (*nonblockingController)(nil)

func NewNonBlockingController(
	index internal.Index,
	aggregator internal.Aggregator,
) *nonblockingController {
	return &nonblockingController{
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
