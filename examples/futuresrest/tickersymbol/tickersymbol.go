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
	fmt.Printf("Symbol: %s", resp.Result.Data.Symbol)
	if resp.Result.Data.Last != nil {
		fmt.Printf(", Last Price: %s", resp.Result.Data.Last)
	}
	if resp.Result.Data.MarkPrice != nil {
		fmt.Printf(", Mark Price: %s", resp.Result.Data.MarkPrice)
	}
	if resp.Result.Data.Bid != nil {
		fmt.Printf(", Bid: %s", resp.Result.Data.Bid)
	}
	if resp.Result.Data.Ask != nil {
		fmt.Printf(", Ask: %s", resp.Result.Data.Ask)
	}
	if resp.Result.Data.IndexPrice != nil {
		fmt.Printf(", Index: %s", resp.Result.Data.IndexPrice)
	}
	if resp.Result.Data.FundingRate != nil {
		fmt.Printf(", Funding Rate: %s", resp.Result.Data.FundingRate)
	}
	if resp.Result.Data.FundingRatePrediction != nil {
		fmt.Printf(", Next Funding Rate: %s", resp.Result.Data.FundingRatePrediction)
	}
	fmt.Printf("\n")
}
