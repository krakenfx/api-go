package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/pkg/derivatives"
)

func main() {
	client := derivatives.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_FUTURES_REST_URL")
	resp, err := client.Tickers()
	if err != nil {
		panic(err)
	}
	for _, ticker := range resp.Tickers {
		fmt.Printf("Symbol: %s, Bid: %s, Ask: %s, Mark: %s, Index: %s\n", ticker.Symbol, ticker.Bid, ticker.Ask, ticker.MarkPrice, ticker.IndexPrice)
	}
}
