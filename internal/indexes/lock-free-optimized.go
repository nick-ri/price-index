package indexes

import (
	"errors"
	"log"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/NickRI/price-index/internal"
	"github.com/NickRI/price-index/internal/indexes/common"
	"github.com/NickRI/price-index/internal/models"
	"github.com/shopspring/decimal"
)

// made just for fun, about 1.5x faster than mb-optimized
type lfOptimized struct {
	ticker models.Ticker

	ringSums *lfRingData

	precision  int
	ratePerSec decimal.Decimal
}

var _ internal.Index = (*lfOptimized)(nil)

func NewLFOptimized(ticker models.Ticker, precision int, ratePerSec int, historyDuration time.Duration) *lfOptimized {
	historySize := int(historyDuration / time.Second)
	return &lfOptimized{
		ticker:     ticker,
		ratePerSec: decimal.NewFromInt(int64(ratePerSec)),
		precision:  precision,
		ringSums:   newLFRingData(historySize, ratePerSec),
	}
}

func (s *lfOptimized) GetTicker() models.Ticker {
	return s.ticker
}

func (s *lfOptimized) SetPrice(tp models.TickerPrice) {
	if tp.Ticker != s.ticker {
		log.Printf("incorrect ticker price readed, got %s, want %s", tp.Ticker, s.ticker)
		return
	}

	s.ringSums.Set(
		tp.Time.Unix(),
		common.ZeroAllocStringToInt(tp.Price, s.precision),
	)
}

func (s *lfOptimized) GetPrice(rStart, rEnd time.Time) (decimal.Decimal, decimal.Decimal, error) {
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

type lfRingData struct {
	size  int
	vsize int
	data  []*data
}

func newLFRingData(size, vsize int) *lfRingData {
	return &lfRingData{
		size:  size,
		vsize: vsize,
		data:  make([]*data, size),
	}
}

func (rd *lfRingData) Set(utime int64, val int64) {
	index := int(utime) % rd.size

	dta := &data{
		val: val,
	}

	for {
		el := (*data)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&rd.data[index]))))

		dta.next = el

		if atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&rd.data[index])),
			unsafe.Pointer(el),
			unsafe.Pointer(dta),
		) {
			break
		}
	}
}

func (rd *lfRingData) Get(utime int64) (sum int64, count int) {
	index := int(utime) % rd.size

	el := (*data)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&rd.data[index]))))
	prev := el

	for el != nil {
		if count == rd.vsize {
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prev.next)), nil)
			break
		}

		sum += el.val
		count++

		prev = el
		el = el.next
	}

	return
}
