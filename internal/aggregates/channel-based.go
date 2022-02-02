package aggregates

import (
	"context"
	"github.com/NickRI/btc_index/internal"
)

type channelBased struct {
	index internal.Index
}

func NewChannelBased(idx internal.Index) *channelBased {
	return &channelBased{
		index: idx,
	}
}

var _ internal.Aggregator = (*channelBased)(nil)

func (s *channelBased) ListenStream(ctx context.Context, subs ...internal.PriceStreamSubscriber) error {
	stream := internal.TickerMultiplexor(ctx, s.index.GetTicker(), subs...)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case tp := <-stream:
			s.index.SetPrice(tp)
		}
	}
}
