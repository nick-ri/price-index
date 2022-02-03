package indexes

import (
	"log"
	"time"

	"github.com/NickRI/price-index/internal"
	"github.com/NickRI/price-index/internal/indexes/common"
	"github.com/NickRI/price-index/internal/models"
	"github.com/shopspring/decimal"
)

type crossRate struct {
	original internal.Index

	crossList []common.Cross

	trustLevel decimal.Decimal // will use cross-rates if original fairness level bellow
}

var _ internal.Index = (*crossRate)(nil)

func NewCrossRate(original internal.Index, trustLevel float64, idxList common.IndexList) *crossRate {
	return &crossRate{
		original:   original,
		crossList:  idxList.SearchCross(original.GetTicker()),
		trustLevel: decimal.NewFromFloat(trustLevel),
	}
}

func (c *crossRate) GetTicker() models.Ticker {
	return c.original.GetTicker()
}

func (c *crossRate) SetPrice(tp models.TickerPrice) {
	c.original.SetPrice(tp)
}

func (c *crossRate) GetPrice(rStart, rEnd time.Time) (decimal.Decimal, decimal.Decimal, error) {
	price, fairness, err := c.original.GetPrice(rStart, rEnd)
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	if fairness.GreaterThanOrEqual(c.trustLevel) {
		return price, fairness, nil
	}

	log.Printf("fairness %s is not enought, try to use cross pairs", fairness)

	for _, cross := range c.crossList {
		price, fairness, err = cross.GetPrice(rStart, rEnd) // TODO: Range by best fair
		if err != nil {
			log.Printf("%s index: cross return price error: %v", c.original.GetTicker(), err)
			continue
		}

		if fairness.LessThan(c.trustLevel) {
			log.Printf("%s index: cross fairness level low: %s < %s", c.original.GetTicker(), fairness, c.trustLevel)
			continue
		}
	}

	return price, fairness, nil // Return prices anyway
}
