package sources

import (
	"errors"
	"github.com/NickRI/btc_index/internal"
	"github.com/NickRI/btc_index/internal/models"
	"github.com/shopspring/decimal"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type sim struct {
	ticker    models.Ticker
	basePrice float64
	errChance float64

	minDelay time.Duration
	maxDelay time.Duration
}

var _ internal.PriceStreamSubscriber = (*sim)(nil)

func NewSim(ticker models.Ticker, basePrice, errChance float64, minDelay, maxDelay time.Duration) *sim {
	return &sim{
		ticker:    ticker,
		basePrice: basePrice,
		errChance: errChance,
		minDelay:  minDelay,
		maxDelay:  maxDelay,
	}
}

func (s *sim) WritePrices(index internal.Index) error {
	go func() {
		for {
			index.SetPrice(models.TickerPrice{
				Ticker: s.ticker,
				Time:   time.Now(),
				Price:  decimal.NewFromFloat(randPrice(s.basePrice)).String(),
			})

			delay := time.Duration(rand.Int63n(int64(s.maxDelay)))

			if delay < s.minDelay {
				delay = s.minDelay
			}

			time.Sleep(delay)
		}
	}()

	return nil
}

func (s *sim) SubscribePriceStream(ticker models.Ticker) (chan models.TickerPrice, chan error) {
	tpCh := make(chan models.TickerPrice)
	errCh := make(chan error)
	go func() {
		for {
			if s.errChance > rand.Float64() {
				errCh <- errors.New("some_error")
				close(tpCh)
				return
			}

			tpCh <- models.TickerPrice{
				Ticker: ticker,
				Time:   time.Now(),
				Price:  decimal.NewFromFloat(randPrice(s.basePrice)).String(),
			}

			delay := time.Duration(rand.Int63n(int64(s.maxDelay)))

			if delay < s.minDelay {
				delay = s.minDelay
			}

			time.Sleep(delay)
		}
	}()

	return tpCh, errCh
}

func randPrice(base float64) float64 {
	k := float64(rand.Intn(5)) * 0.01
	r := rand.Intn(1)
	if r == 2 {
		k = -k
	}

	pow := math.Pow(10, float64(models.FiatDecimals))
	intermed := (base * (1 + k)) * pow

	return math.Ceil(intermed) / pow
}
