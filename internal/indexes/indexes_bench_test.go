package indexes

import (
	"sync"
	"testing"
	"time"

	"github.com/NickRI/price-index/internal"
	"github.com/NickRI/price-index/internal/models"
	"github.com/shopspring/decimal"
)

const amountOfStreams = 100

var one = decimal.NewFromInt(1)

func Benchmark_mutexBased_GetPrice(b *testing.B) {
	mbased := NewMutexBased(models.Ticker(0), amountOfStreams, time.Duration(b.N)*time.Second)
	getPriceBench(b, mbased)
}

func Benchmark_mutexBased_SetPrice(b *testing.B) {
	mbased := NewMutexBased(models.Ticker(0), amountOfStreams, time.Duration(b.N)*time.Second)
	setPriceBench(b, mbased)
}

func Benchmark_mbOptimized_GetPrice(b *testing.B) {
	mbased := NewMBOptimized(models.Ticker(0), models.FiatDecimals, amountOfStreams, time.Duration(b.N)*time.Second)
	getPriceBench(b, mbased)
}

func Benchmark_mbOptimized_SetPrice(b *testing.B) {
	mbased := NewMBOptimized(models.Ticker(0), models.FiatDecimals, amountOfStreams, time.Duration(b.N)*time.Second)
	setPriceBench(b, mbased)
}

func Benchmark_lfOptimized_GetPrice(b *testing.B) {
	mbased := NewLFOptimized(models.Ticker(0), models.FiatDecimals, amountOfStreams, time.Duration(b.N)*time.Second)
	getPriceBench(b, mbased)
}

func Benchmark_lfOptimized_SetPrice(b *testing.B) {
	mbased := NewLFOptimized(models.Ticker(0), models.FiatDecimals, amountOfStreams, time.Duration(b.N)*time.Second)
	setPriceBench(b, mbased)
}

func setPriceBench(b *testing.B, index internal.Index) {
	b.ReportAllocs()
	b.SetParallelism(amountOfStreams)

	prebuit := decimal.NewFromFloat(37948.12).String()

	zeroTime := time.Unix(0, 0)

	wg := sync.WaitGroup{}
	wg.Add(amountOfStreams)

	for i := 0; i < amountOfStreams; i++ {
		go func() {
			for sec := 0; sec < b.N; sec++ {
				index.SetPrice(models.TickerPrice{
					Ticker: models.Ticker(0),
					Price:  prebuit,
					Time:   zeroTime.Add(time.Duration(sec) * time.Second),
				})
			}

			wg.Done()
		}()
	}

	wg.Wait()
}

func getPriceBench(b *testing.B, index internal.Index) {
	b.ReportAllocs()

	zeroTime := time.Unix(0, 0)

	prebuit := decimal.NewFromFloat(37948.14).String()

	for sec := 0; sec < b.N; sec++ {
		for i := 0; i < amountOfStreams; i++ {
			index.SetPrice(models.TickerPrice{
				Ticker: models.Ticker(0),
				Price:  prebuit,
				Time:   zeroTime.Add(time.Duration(sec) * time.Second),
			})
		}
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			price, fair, err := index.GetPrice(zeroTime, zeroTime.Add(time.Duration(b.N)*time.Second))
			if err != nil {
				b.Fatal(err)
			}

			if price.IsZero() {
				b.Fatal("price can't be zero")
			}

			if fair.LessThan(one) {
				b.Fatal("fairness can't be less than 1")
			}
		}
	})
}
