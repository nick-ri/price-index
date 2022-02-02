package internal

import (
	"context"
	"github.com/NickRI/btc_index/internal/models"
	"github.com/shopspring/decimal"
	"time"
)

// PriceStreamSubscriber base stream subscribe method
type PriceStreamSubscriber interface {
	SubscribePriceStream(models.Ticker) (chan models.TickerPrice, chan error)
}

// PirceWriter for more efficient data setting solution
type PirceWriter interface {
	WritePrices(index Index) error
}

// Aggregator provides ability to listen multipe streams
type Aggregator interface {
	ListenStream(ctx context.Context) error
}

// Index main element to store prices and calculate average price
type Index interface {
	GetTicker() models.Ticker
	SetPrice(tp models.TickerPrice)
	// Added fairnes to show fitness of current index
	GetPrice(rStart, rEnd time.Time) (price, fairnsess decimal.Decimal, err error)
}

// Controller high-level components that's make aggregation and indexes works together
type Controller interface {
	Exec(ctx context.Context) error
}
