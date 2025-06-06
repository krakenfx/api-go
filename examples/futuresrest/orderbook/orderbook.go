package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/pkg/derivatives"
	"github.com/krakenfx/api-go/pkg/kraken"
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
	for _, side := range [][][]*kraken.Money{resp.OrderBook.Asks, resp.OrderBook.Bids} {
		for i := 9; i >= 0 && i < len(side); i-- {
			fmt.Printf("%s - %s\n", side[i][0], side[i][1])
		}
		fmt.Printf("---\n")
	}
}
