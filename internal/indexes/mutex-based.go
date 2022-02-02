package indexes

import (
	"errors"
	"github.com/NickRI/btc_index/internal"
	"github.com/NickRI/btc_index/internal/models"
	"github.com/shopspring/decimal"
	"log"
	"sync"
	"time"
)

type mutexBased struct {
	mu sync.Mutex

	prices         map[int64]models.Numbers
	ticker         models.Ticker
	historyTTLSecs int64
	sourcesCount   decimal.Decimal
}

var _ internal.Index = (*mutexBased)(nil)

func NewMutexBased(ticker models.Ticker, sourcesCount int64, historyDuration time.Duration) *mutexBased {
	return &mutexBased{
		ticker:         ticker,
		sourcesCount:   decimal.NewFromInt(sourcesCount),
		historyTTLSecs: int64(historyDuration / time.Second),
		prices:         make(map[int64]models.Numbers),
	}
}

func (s *mutexBased) GetTicker() models.Ticker {
	return s.ticker
}

func (s *mutexBased) SetPrice(tp models.TickerPrice) {
	if tp.Ticker != s.ticker {
		log.Printf("incorrect ticker price readed, got %s, want %s", tp.Ticker, s.ticker)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	utime := tp.Time.Unix()

	s.prices[utime] = append(s.prices[utime], decimal.RequireFromString(tp.Price))

	if diff := len(s.prices) - int(s.historyTTLSecs); diff > 0 {
		for i := 0; i < diff; i++ {
			delete(s.prices, utime-s.historyTTLSecs-int64(i)) // cleanup history
		}
	}
}

func (s *mutexBased) GetPrice(rStart, rEnd time.Time) (decimal.Decimal, decimal.Decimal, error) {
	uStart := rStart.Unix()
	uEnd := rEnd.Unix()

	if uStart == uEnd {
		return decimal.Zero, decimal.Zero, errors.New("empty range provided")
	}

	priceResult := models.PriceResult{
		Prices:   make(models.Numbers, 0, uEnd-uStart),
		Fairness: make(models.Numbers, 0, uEnd-uStart),
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for usec := uStart; usec < uEnd; usec++ {
		fairness := decimal.NewFromInt(int64(len(s.prices[usec]))).Div(s.sourcesCount)
		secPrice := s.prices[usec].Avg()

		if !secPrice.IsZero() {
			priceResult.Prices = append(priceResult.Prices, secPrice)
			priceResult.Fairness = append(priceResult.Fairness, fairness)
		}
	}

	return priceResult.Prices.Avg(), priceResult.Fairness.Avg(), nil
}
