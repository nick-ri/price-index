package aggregates

import (
	"context"

	"github.com/NickRI/btc_index/internal"
)

type efficient struct {
	index   internal.Index
	writers []internal.PirceWriter
}

func NewEfficient(idx internal.Index, writers ...internal.PirceWriter) *efficient {
	return &efficient{
		index:   idx,
		writers: writers,
	}
}

var _ internal.Aggregator = (*efficient)(nil)

func (s *efficient) ListenStream(ctx context.Context) error {
	for _, writer := range s.writers {
		if err := writer.WritePrices(s.index); err != nil {
			return err
		}
	}

	return nil
}
