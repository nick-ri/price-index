package common

import (
	"errors"
	"fmt"
	"time"

	"github.com/NickRI/price-index/internal"
	"github.com/NickRI/price-index/internal/models"
	"github.com/shopspring/decimal"
)

type CrossType int8

const (
	Reversed CrossType = iota + 1
	Direct
	Mixed
)

type Cross struct {
	Type          CrossType
	First, Second IndexDir
}

type IndexDir struct {
	internal.Index
	NeedReverse bool
}

func (c *Cross) GetPrice(rStart, rEnd time.Time) (price, fairness decimal.Decimal, err error) {
	fprice, ffairness, err := c.First.GetPrice(rStart, rEnd)
	if err != nil {
		return decimal.Zero, decimal.Zero, fmt.Errorf("first index return error: %w", err)
	}

	if fprice.IsZero() {
		return decimal.Zero, decimal.Zero, errors.New("price is zero")
	}

	if c.First.NeedReverse {
		fprice = decimal.NewFromInt(1).Div(fprice)
	}

	sprice, sfairness, err := c.Second.GetPrice(rStart, rEnd)
	if err != nil {
		return decimal.Zero, decimal.Zero, fmt.Errorf("second index return error: %w", err)
	}

	if sprice.IsZero() {
		return decimal.Zero, decimal.Zero, errors.New("price is zero")
	}

	if c.Second.NeedReverse {
		sprice = decimal.NewFromInt(1).Div(sprice)
	}

	switch c.Type {
	case Reversed, Direct:
		price = fprice.Div(sprice)
	case Mixed:
		price = fprice.Mul(sprice)
	}

	if ffairness.GreaterThan(sfairness) { // need to get smallest
		fairness = sfairness
	} else {
		fairness = ffairness
	}

	return
}

// BTC/USD = ETH/USD / ETH/BTC  reversed
// BTC/USD = BTC/ETH / USD/ETH  directed
// BTC/USD = BTC/ETH * ETH/USD  mixed

func (l IndexList) SearchCross(orig models.Ticker) []Cross {
	crossList := make([]Cross, 0)

	for _, fidx := range l { // N^2
		for _, sidx := range l {
			switch {
			case orig.Secondary() == fidx.Secondary &&
				orig.Base() == sidx.Secondary &&
				fidx.Base == sidx.Base:

				crossList = append(crossList, Cross{
					Type:   Reversed,
					First:  IndexDir{Index: fidx, NeedReverse: fidx.Reverted},
					Second: IndexDir{Index: sidx, NeedReverse: sidx.Reverted},
				})

			case orig.Base() == fidx.Base &&
				orig.Secondary() == sidx.Base &&
				fidx.Secondary == sidx.Secondary:

				crossList = append(crossList, Cross{
					Type:   Direct,
					First:  IndexDir{Index: fidx, NeedReverse: fidx.Reverted},
					Second: IndexDir{Index: sidx, NeedReverse: sidx.Reverted},
				})

			case orig.Base() == fidx.Base &&
				orig.Secondary() == sidx.Secondary &&
				fidx.Secondary == sidx.Base:

				crossList = append(crossList, Cross{
					Type:   Mixed,
					First:  IndexDir{Index: fidx, NeedReverse: fidx.Reverted},
					Second: IndexDir{Index: sidx, NeedReverse: sidx.Reverted},
				})
			}
		}
	}

	set := make(map[models.Ticker]struct{}) // need to dedup for range try
	resultList := make([]Cross, 0, len(crossList))

	for _, item := range crossList {
		_, fexists := set[item.First.Index.GetTicker()]
		_, sexists := set[item.Second.Index.GetTicker()]

		if fexists && sexists {
			continue
		}

		set[item.First.Index.GetTicker()] = struct{}{}
		set[item.Second.Index.GetTicker()] = struct{}{}

		resultList = append(resultList, item)
	}

	return resultList
}
