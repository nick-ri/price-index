package indexes

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/NickRI/btc_index/internal"
	"github.com/NickRI/btc_index/internal/indexes/common"
	"github.com/NickRI/btc_index/internal/models"
	"github.com/shopspring/decimal"
)

// mutex-based index with optimizations, about 5x times faster on my laptop
type mbOptimized struct {
	ticker models.Ticker

	ringSums *ringData

	precision  int
	ratePerSec decimal.Decimal
}

var _ internal.Index = (*mbOptimized)(nil)

func NewMBOptimized(ticker models.Ticker, precision int, ratePerSec int, historyDuration time.Duration) *mbOptimized {
	historySize := int(historyDuration / time.Second)

	return &mbOptimized{
		ticker:     ticker,
		ratePerSec: decimal.NewFromInt(int64(ratePerSec)),
		precision:  precision,
		ringSums:   newRingData(historySize, ratePerSec),
	}
}

func (s *mbOptimized) GetTicker() models.Ticker {
	return s.ticker
}

func (s *mbOptimized) SetPrice(tp models.TickerPrice) {
	if tp.Ticker != s.ticker {
		log.Printf("incorrect ticker price readed, got %s, want %s", tp.Ticker, s.ticker)
		return
	}

	s.ringSums.Set(
		tp.Time.Unix(),
		common.ZeroAllocStringToInt(tp.Price, s.precision),
	)
}

func (s *mbOptimized) GetPrice(rStart, rEnd time.Time) (decimal.Decimal, decimal.Decimal, error) {
	uStart := rStart.Unix()
	uEnd := rEnd.Unix()

	// TODO: Check range possible in terms of the ring loop

	if uStart == uEnd {
		return decimal.Zero, decimal.Zero, errors.New("empty range provided")
	}

	var (
		priceSum  int64
		countsSum int
	)

	for usec := uStart; usec < uEnd; usec++ {
		sum, count := s.ringSums.Get(usec)
		priceSum += sum
		countsSum += count
	}

	if priceSum == 0 {
		return decimal.Zero, decimal.Zero, nil
	}

	length := decimal.NewFromInt(uEnd - uStart)

	count := decimal.NewFromInt(int64(countsSum))

	return decimal.New(priceSum, -int32(s.precision)).Div(count),
		count.Div(length.Mul(s.ratePerSec)),
		nil
}

const muSize = 8 // the size is enough. For random r/w it should shows best results

type ringData struct {
	mu    [muSize]sync.Mutex
	size  int
	vsize int
	data  []*data
}

type data struct {
	val  int64
	next *data
}

func newRingData(size, vsize int) *ringData {
	return &ringData{
		size:  size,
		vsize: vsize,
		data:  make([]*data, size),
	}
}

func (rd *ringData) Set(utime, val int64) {
	index := int(utime) % rd.size

	rd.mu[utime%muSize].Lock()

	el := rd.data[index]

	rd.data[index] = &data{
		val:  val,
		next: el,
	}

	rd.mu[utime%muSize].Unlock()
}

func (rd *ringData) Get(utime int64) (sum int64, count int) {
	index := int(utime) % rd.size

	rd.mu[utime%muSize].Lock()

	el := rd.data[index]
	prev := el

	for el != nil {
		if count == rd.vsize {
			prev.next = nil
			break // stop
		}

		sum += el.val
		count++

		prev = el
		el = el.next
	}

	rd.mu[utime%muSize].Unlock()

	return
}
