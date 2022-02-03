package aggregates

import (
	"context"

	"github.com/NickRI/price-index/internal"
)

type channelBased struct {
	index internal.Index
	subs  []internal.PriceStreamSubscriber
}

func NewChannelBased(idx internal.Index, subs ...internal.PriceStreamSubscriber) *channelBased {
	return &channelBased{
		index: idx,
		subs:  subs,
	}
}

var _ internal.Aggregator = (*channelBased)(nil)

func (s *channelBased) ListenStream(ctx context.Context) error {
	stream := internal.TickerMultiplexor(ctx, s.index.GetTicker(), s.subs...)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case tp := <-stream:
			s.index.SetPrice(tp)
		}
	}
}
