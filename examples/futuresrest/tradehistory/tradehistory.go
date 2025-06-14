package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/pkg/derivatives"
)

// Derivative contract.
var contract = "PF_XBTUSD"

func main() {
	client := derivatives.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_FUTURES_REST_URL")
	fmt.Printf("> Fetch trades history.\n")
	resp, err := client.TradeHistory(&derivatives.TradeHistoryRequest{
		Symbol: contract,
	})
	if err != nil {
		panic(err)
	}
	for _, trade := range resp.Result.History {
		fmt.Printf("%s %s units of %s on %s at price %s\n", trade.Side, trade.Side, contract, trade.Time, trade.Price)
	}
}
