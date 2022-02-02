package common

import (
	"github.com/NickRI/btc_index/internal"
)

type IndexItem struct {
	Base      string
	Secondary string
	Reverted  bool
	internal.Index
}

type IndexList []IndexItem

func (l *IndexList) Add(idx internal.Index) {
	ticker := idx.GetTicker()

	*l = append(*l, IndexItem{
		Base:      ticker.Base(),
		Secondary: ticker.Secondary(),
		Index:     idx,
	})

	*l = append(*l, IndexItem{
		Base:      ticker.Secondary(),
		Secondary: ticker.Base(),
		Reverted:  true,
		Index:     idx,
	})
}

func (l *IndexList) AddItems(idxs ...internal.Index) {
	for _, idx := range idxs {
		l.Add(idx)
	}
}
