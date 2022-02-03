package internal

import (
	"context"
	"log"
	"time"

	"github.com/shopspring/decimal"
)

type Result struct {
	Price    decimal.Decimal
	Fairness decimal.Decimal
}

func RangePriceEmitter(ctx context.Context, index Index, ranging time.Duration) chan Result {
	out := make(chan Result)

	go func() {
		defer close(out)

		now := time.Now()

		if err := sleepCtx(ctx, now.Sub(now.Truncate(ranging))); err != nil { // need to align ticker start
			log.Printf("sleepCtx: error occurred, exiting : %v", err)
			return
		}

		ticker := time.NewTicker(ranging)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				now = time.Now()
				price, fair, err := index.GetPrice(now.Add(-ranging), now)
				if err != nil {
					log.Printf("index.GetPrice: error occurred, exiting : %v", err)
					return
				}

				if price.IsZero() { // for small ranges
					continue
				}

				out <- Result{
					Price:    price,
					Fairness: fair,
				}
			}
		}
	}()

	return out
}

func sleepCtx(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
