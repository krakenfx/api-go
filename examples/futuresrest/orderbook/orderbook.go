package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/v2/pkg/decimal"
	"github.com/krakenfx/api-go/v2/pkg/derivatives"
)

// Derivative contract.
var contract string = "PF_XBTUSD"

func main() {
	client := derivatives.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_FUTURES_REST_URL")
	resp, err := client.OrderBook(&derivatives.OrderBookRequest{
		Symbol: contract,
	})
	if err != nil {
		panic(err)
	}
	for _, side := range [][][]*decimal.Decimal{resp.Result.OrderBook.Asks, resp.Result.OrderBook.Bids} {
		for i := 9; i >= 0 && i < len(side); i-- {
			fmt.Printf("%s - %s\n", side[i][0], side[i][1])
		}
		fmt.Printf("---\n")
	}
}
