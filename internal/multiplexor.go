package internal

import (
	"context"
	"github.com/NickRI/btc_index/internal/models"
	"log"
)

func TickerMultiplexor(ctx context.Context, ticker models.Ticker, subscibers ...PriceStreamSubscriber) chan models.TickerPrice {
	output := make(chan models.TickerPrice, 1)

	for _, sub := range subscibers {
		go func(tpCh chan models.TickerPrice, errCh chan error) {
			for {
				select {
				case <-ctx.Done():
					return

				case tp, ok := <-tpCh:
					if !ok {
						log.Println("multiplexor: accidentialy closed stream channel, finalizing")
						return
					}

					output <- tp
				case err := <-errCh:
					log.Printf("multiplexor: error occurred, finalizing: %v", err)
					return
				}
			}
		}(sub.SubscribePriceStream(ticker))
	}

	return output
}
