package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/pkg/derivatives"
)

// Derivative contract.
var contract string = "PF_XBTUSD"

func main() {
	client := derivatives.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_FUTURES_REST_URL")
	resp, err := client.TickerSymbol(contract)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Symbol: %s", resp.Ticker.Symbol)
	if resp.Ticker.Last != nil {
		fmt.Printf(", Last Price: %s", resp.Ticker.Last)
	}
	if resp.Ticker.MarkPrice != nil {
		fmt.Printf(", Mark Price: %s", resp.Ticker.MarkPrice)
	}
	if resp.Ticker.Bid != nil {
		fmt.Printf(", Bid: %s", resp.Ticker.Bid)
	}
	if resp.Ticker.Ask != nil {
		fmt.Printf(", Ask: %s", resp.Ticker.Ask)
	}
	if resp.Ticker.IndexPrice != nil {
		fmt.Printf(", Index: %s", resp.Ticker.IndexPrice)
	}
	if resp.Ticker.FundingRate != nil {
		fmt.Printf(", Funding Rate: %s", resp.Ticker.FundingRate)
	}
	if resp.Ticker.FundingRatePrediction != nil {
		fmt.Printf(", Next Funding Rate: %s", resp.Ticker.FundingRatePrediction)
	}
	fmt.Printf("\n")
}
